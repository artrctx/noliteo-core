-- Will be using built in gen_random_uuid()
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION auto_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql'; 


CREATE TABLE usr(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Hashed Token
    token VarChar(32),
    username Text,
    created_at timestamptz NOT NULL DEFAULT now()
);