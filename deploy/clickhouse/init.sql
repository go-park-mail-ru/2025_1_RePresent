CREATE TABLE reviews (
    id UUID DEFAULT generateUUIDv4(),
    question text,
    page text,
    comment VARCHAR(200),
    rating INT,
    created_at TIMESTAMP DEFAULT toTimeZone(now(), 'Europe/Moscow')
) ENGINE = MergeTree()
ORDER BY (created_at);
