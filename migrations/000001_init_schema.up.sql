-- Will be using built in gen_random_uuid()
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION public.auto_updated_at () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE token (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Hashed Token
    key VarChar(60) UNIQUE,
    ident Text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION public.validate_token_key (tkn text) RETURNS TABLE (id uuid, ident text) LANGUAGE plpgsql AS $function$
BEGIN
    RETURN QUERY SELECT t.id, t.ident FROM "token" t WHERE t.key=crypt(tkn, t.key) LIMIT 1;
END;

$function$;

CREATE OR REPLACE FUNCTION public.hash_token_trigger () RETURNS TRIGGER LANGUAGE plpgsql AS $function$
BEGIN

IF TG_OP = 'UPDATE' THEN
    RAISE EXCEPTION 'token key cannot be updated';
END IF;

-- THROW IF DUPLICATE KEY EXISTS
IF EXISTS ( SELECT 1 FROM validate_token_key (new."key") ) THEN 
    RAISE EXCEPTION 'duplicate key detected';
END IF;

new."key" := crypt(new.key, gen_salt('bf'));

RETURN new;

END;
$function$;

CREATE OR REPLACE TRIGGER hash_token_trg BEFORE INSERT
OR
UPDATE ON public.token FOR EACH ROW
EXECUTE FUNCTION hash_token_trigger ();