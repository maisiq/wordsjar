-- +goose Up
-- +goose StatementBegin
CREATE TABLE words (
    id UUID PRIMARY KEY,
    en varchar UNIQUE,
    ru VARCHAR[],
    transcription varchar,
    examples TEXT[]
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar UNIQUE,
    password text,
    is_admin bool DEFAULT false
);

CREATE TABLE user_words (
    id bigint GENERATED ALWAYS AS IDENTITY,
    username varchar REFERENCES users(username) ON DELETE CASCADE,
    word_en varchar REFERENCES words(en) ON DELETE CASCADE,
    knowledge_rating real DEFAULT 0.0,
    consecutive_success_attempts smallint DEFAULT 1,
    last_attempt timestamp DEFAULT (now() AT TIME ZONE 'utc' - INTERVAL '1 day 1 hour'),
    UNIQUE (username, word_en)
);

-- +goose StatementEnd

