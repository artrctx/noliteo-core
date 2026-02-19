CREATE TYPE rtc_type AS ENUM('offer', 'answer');

CREATE TABLE rtc_description (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sdp TEXT NOT NULL,
    type rtc_type NOT NULL,
    token_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- eastablish token relation
ALTER TABLE rtc_description
ADD CONSTRAINT fk_token FOREIGN KEY (token_id) REFERENCES token (id) ON DELETE CASCADE;