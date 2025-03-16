package repository

import (
	"database/sql"
	"strings"
	"time"

	"RepoTracker/src/entity"
)

type RepoStatsRepository struct {
	db *sql.DB
}

func NewRepoStatsRepository(db *sql.DB) *RepoStatsRepository {
	return &RepoStatsRepository{db: db}
}

func (r *RepoStatsRepository) GetByMint(repoFullName string) ([]entity.RepoStats, error) {
	query := `
		SELECT id, repo_name, repo_owner, repo_full_name, stars_count, forks_count, 
		       contributors_count, stats_date, updated_at
		FROM repo_stats
		WHERE repo_full_name = $1
		ORDER BY stats_date DESC
	`

	rows, err := r.db.Query(query, repoFullName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []entity.RepoStats
	for rows.Next() {
		var s entity.RepoStats
		var statsDate, updatedAt time.Time

		err := rows.Scan(
			&s.ID, &s.RepoName, &s.RepoOwner, &s.RepoFullName,
			&s.StarsCount, &s.ForksCount, &s.ContributorsCount,
			&statsDate, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		s.StatsDate = statsDate
		s.UpdatedAt = updatedAt
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *RepoStatsRepository) GetAll() ([]entity.RepoStats, error) {
	query := `
		SELECT id, repo_name, repo_owner, repo_full_name, stars_count, forks_count, 
		       contributors_count, stats_date, updated_at
		FROM repo_stats
		ORDER BY repo_full_name, stats_date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []entity.RepoStats
	for rows.Next() {
		var s entity.RepoStats
		var statsDate, updatedAt time.Time

		err := rows.Scan(
			&s.ID, &s.RepoName, &s.RepoOwner, &s.RepoFullName,
			&s.StarsCount, &s.ForksCount, &s.ContributorsCount,
			&statsDate, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		s.StatsDate = statsDate
		s.UpdatedAt = updatedAt
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *RepoStatsRepository) SaveStats(repoURL string, dataPoint entity.RepoDataPoint) error {
	parts := strings.Split(repoURL, "/")
	var repoOwner, repoName string

	if len(parts) >= 2 {
		repoName = parts[len(parts)-1]
		repoOwner = parts[len(parts)-2]
	} else {
		repoName = "unknown"
		repoOwner = "unknown"
	}

	query := `
		INSERT INTO repo_stats (
			repo_name, repo_owner, stars_count, forks_count, contributors_count, stats_date
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (repo_full_name, stats_date) 
		DO UPDATE SET
			stars_count = EXCLUDED.stars_count,
			forks_count = EXCLUDED.forks_count,
			contributors_count = EXCLUDED.contributors_count,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(
		query,
		repoName,
		repoOwner,
		dataPoint.Stars,
		dataPoint.Forks,
		dataPoint.Contributors,
		dataPoint.Time,
	)

	return err
}
