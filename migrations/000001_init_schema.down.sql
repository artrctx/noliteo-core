DROP TRIGGER IF EXISTS hash_token_trg ON public.token;

DROP FUNCTION IF EXISTS public.hash_token_trigger;

DROP FUNCTION IF EXISTS public.validate_token_key;

DROP TABLE IF EXISTS token;

DROP FUNCTION IF EXISTS public.auto_updated_at;

DROP EXTENSION IF EXISTS pgcrypto;
