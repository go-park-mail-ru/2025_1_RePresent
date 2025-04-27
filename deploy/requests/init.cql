CREATE KEYSPACE IF NOT EXISTS link_space WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE link_space;

CREATE TABLE IF NOT EXISTS links (
    link text PRIMARY KEY,
    user_id int,
    height int,
    width int
);


CREATE INDEX IF NOT EXISTS user_id_index ON links (user_id);
