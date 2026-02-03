
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
    username varchar REFERENCES users(username),
    word_en varchar REFERENCES words(en),
    knowledge_rating real DEFAULT 0.0,
    consecutive_success_attempts smallint DEFAULT 1,
    last_attempt timestamp DEFAULT (now() AT TIME ZONE 'utc' - INTERVAL '1 day 1 hour'),
    UNIQUE (username, word_en)
);

-- TEST DATA

-- USERS

INSERT INTO users(username, password, is_admin)
VALUES
-- REGULAR USERS
('testuser', 'secret', false),
-- ADMINS
('admin', 'secret', true);


-- WORDS

INSERT INTO words
VALUES ('73af303e-63c6-4419-8f19-158dbe1f2a3c', 'word', '{слово}', 'smth', null);  -- Add word with known id

INSERT INTO words
SELECT gen_random_uuid(), 'word' || t.id::text, ('{слово' || t.id::text || '}')::text[], 'smth' || t.id::text, null
FROM generate_series(1, 30) t(id);


-- USER WORDS

INSERT INTO user_words(username, word_en, last_attempt)
VALUES
-- Word included in TestMode
('testuser', 'word', DEFAULT),
('testuser', 'word1', DEFAULT),
 -- Word not included in TestMode
('testuser', 'word2', now() AT TIME ZONE 'utc');
