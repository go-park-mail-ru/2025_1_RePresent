CREATE KEYSPACE IF NOT EXISTS slot_space 
WITH replication = {
    'class': 'SimpleStrategy', 
    'replication_factor': 1
};

USE slot_space;

-- Таблица слотов
CREATE TABLE IF NOT EXISTS slots (
    link uuid PRIMARY KEY,
    user_id int NOT NULL,
    format_code uuid NOT NULL,
    min_price decimal(12,2) NOT NULL,
    created_at timestamp
);

-- WITH compression = {
--    'sstable_compression': 'LZ4Compressor',
--    'chunk_length_kb': 32
--};

-- Таблица форматов
CREATE TABLE IF NOT EXISTS formats (
    code uuid PRIMARY KEY,
    height int NOT NULL,
    width int NOT NULL,
    description text
);

-- Материальное представление по пользователям
CREATE MATERIALIZED VIEW IF NOT EXISTS slots_by_user_time AS
    SELECT link, user_id, format_code, min_price, created_at 
    FROM slots
    WHERE user_id IS NOT NULL AND link IS NOT NULL
    PRIMARY KEY (user_id, created_at, link)
    WITH CLUSTERING ORDER BY (created_at DESC);

