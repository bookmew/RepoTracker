CREATE OR REPLACE VIEW repo_stats_latest AS
SELECT *
FROM repo_stats
WHERE stats_date = CURRENT_DATE;