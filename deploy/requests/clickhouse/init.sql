CREATE DATABASE IF NOT EXISTS csat;

CREATE TABLE IF NOT EXISTS csat.reviews (
    id UUID DEFAULT generateUUIDv4(),
    user_id INT,
    question text,
    page text,
    comment VARCHAR(200),
    rating INT,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (created_at);
