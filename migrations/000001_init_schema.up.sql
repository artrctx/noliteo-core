-- Will be using built in gen_random_uuid()
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION public.auto_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql'; 

CREATE OR REPLACE FUNCTION public.hash_token_trigger()
RETURN TRIGGER
LANGAUAGE plpgsql
AS $function$DECLARE
    api_key text;
BEGIN

END;
$function$;


CREATE TABLE usr(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Hashed Token
    token VarChar(32),
    username Text,
    created_at timestamptz NOT NULL DEFAULT now()
);


-- ```-- Hashing apikey on insertCREATE OR REPLACE FUNCTION public.hash_apikey_trigger() RETURNS trigger LANGUAGE plpgsqlAS $function$DECLARE api_key text;BEGIN IF new.type='SANDBOX' THEN api_key := 'sk_test_' || new.secret; ELSE api_key := 'sk_live_' || new.secret; END IF;
--  IF EXISTS (SELECT * FROM validate_apikey_hash(('{"key": "' || api_key || '"}')::jsonb)) THEN RAISE EXCEPTION 'Duplicate active API key detected'; END IF;
--  -- Add prefix as identifier new.key_prefix := substring(api_key, 1, 7) || '...' || substring(api_key, LENGTH(api_key) - 5);
--  -- Generate a bcrypt hash new.secret := crypt(api_key, gen_salt('bf'));
--  return new;END;
-- ``````CREATE OR REPLACE FUNCTION public.validate_apikey_hash(args jsonb) RETURNS TABLE(id text, org_id text, type environment) LANGUAGE plpgsqlAS $function$DECLARE hash text;BEGIN -- https://www.postgresql.org/docs/9.5/functions-json.html RETURN QUERY SELECT a.id, a.org_id, a.type FROM "ApiKeys" a WHERE a.enabled = true AND (expiration > NOW() OR expiration IS NULL) AND (args#>>'{orgIds,0}' IS NULL OR (args->'orgIds') ? a.org_id) AND (a.secret = crypt(args->>'key', a.secret));END;
-- $function$;```