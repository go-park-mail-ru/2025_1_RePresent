CREATE MATERIALIZED VIEW IF NOT EXISTS links_by_user_id AS
SELECT user_id, link
FROM links
WHERE user_id IS NOT NULL AND link IS NOT NULL
PRIMARY KEY (user_id, link);