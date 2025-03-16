CREATE TABLE repo_stats (
    id SERIAL PRIMARY KEY,
    repo_name TEXT NOT NULL,
    repo_owner TEXT NOT NULL,
    repo_full_name TEXT GENERATED ALWAYS AS (repo_owner || '/' || repo_name) STORED,
    stars_count BIGINT DEFAULT 0,
    forks_count BIGINT DEFAULT 0,
    contributors_count BIGINT DEFAULT 0,
    stats_date DATE NOT NULL DEFAULT CURRENT_DATE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_repo_stats_repo_id_date ON repo_stats(repo_full_name, stats_date);

CREATE INDEX idx_repo_stats_repo_full_name ON repo_stats(repo_full_name);
CREATE INDEX idx_repo_stats_stats_date ON repo_stats(stats_date);
CREATE INDEX idx_repo_stats_stars_count ON repo_stats(stars_count);
CREATE INDEX idx_repo_stats_forks_count ON repo_stats(forks_count);

-- Insert and update sample data
-- INSERT INTO repo_stats (
--     repo_name, 
--     repo_owner, 
--     stars_count, 
--     forks_count, 
--     contributors_count, 
--     stats_date
-- ) VALUES (
--     'awesome-project', 
--     'github-user', 
--     1250, 
--     320, 
--     150, 
--     CURRENT_DATE
-- )
-- ON CONFLICT (repo_full_name, stats_date) 
-- DO UPDATE SET
--     stars_count = EXCLUDED.stars_count,
--     forks_count = EXCLUDED.forks_count,
--     contributors_count = EXCLUDED.contributors_count;

-- Query sample data
-- SELECT stats_date, stars_count 
-- FROM repo_stats 
-- WHERE repo_id = 123456789
-- ORDER BY stats_date;