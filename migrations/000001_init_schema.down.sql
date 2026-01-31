-- drop auth updated at trigger function
DROP TABLE IF EXISTS usr;

DROP FUNCTION IF EXISTS public.hash_token_trigger;

DROP FUNCTION IF EXISTS public.auto_updated_at;

DROP EXTENSION IF EXISTS pgcrypto;
