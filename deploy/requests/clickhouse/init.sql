CREATE DATABASE IF NOT EXISTS csat;

CREATE TABLE IF NOT EXISTS csat.reviews (
    id UUID DEFAULT generateUUIDv4(),
    user_id INT,
    question text,
    page text,
    comment VARCHAR(200),
    rating INT,
    created_at TIMESTAMP DEFAULT toTimeZone(now(), 'Europe/Moscow')
) ENGINE = MergeTree()
ORDER BY (created_at);


CREATE DATABASE IF NOT EXISTS pay;

CREATE TABLE IF NOT EXISTS csat.activity (
    id UUID DEFAULT generateUUIDv4(),
    created_at TIMESTAMP DEFAULT toTimeZone(now(), 'Europe/Moscow'),
    user_indent text DEFAULT randomString(),
    cash_off DECIMAL(12, 2)
) ENGINE = MergeTree()
ORDER BY created_at;
