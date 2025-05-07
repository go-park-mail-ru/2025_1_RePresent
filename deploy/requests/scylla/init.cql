CREATE KEYSPACE IF NOT EXISTS slot_space 
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

USE slot_space;

CREATE TABLE IF NOT EXISTS slots (
    link uuid PRIMARY KEY,
    user_id int,
    slot_name text,
    format_code smallint,
    min_price decimal,
    is_active boolean,
    created_at timestamp
) WITH compaction = { 'class' : 'LeveledCompactionStrategy' };

CREATE TABLE IF NOT EXISTS user_links (
    user_id int,
    created_at timestamp,
    link uuid,
    PRIMARY KEY (user_id, created_at, link)
) WITH CLUSTERING ORDER BY (created_at DESC)
AND compaction = { 'class' : 'LeveledCompactionStrategy' };

CREATE TABLE IF NOT EXISTS formats (
    code smallint PRIMARY KEY,
    height int,
    width int,
    description text
) WITH compaction = { 'class' : 'LeveledCompactionStrategy' };

-- Предзаполнение таблицы форматов
INSERT INTO formats (code, height, width, description) VALUES (1, 300, 300, 'Квадратный формат');
INSERT INTO formats (code, height, width, description) VALUES (2, 300, 300, 'Прямоугольный формат');
INSERT INTO formats (code, height, width, description) VALUES (3, 300, 300, 'Широкий формат');
INSERT INTO formats (code, height, width, description) VALUES (4, 300, 300, 'Высокий формат');
