package repository

import (
	"database/sql"
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
		       contributors_count, stats_date, created_at
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
		var statsDate, createdAt time.Time

		err := rows.Scan(
			&s.ID, &s.RepoName, &s.RepoOwner, &s.RepoFullName,
			&s.StarsCount, &s.ForksCount, &s.ContributorsCount,
			&statsDate, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		s.StatsDate = statsDate
		s.CreatedAt = createdAt
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
		       contributors_count, stats_date, created_at
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
		var statsDate, createdAt time.Time

		err := rows.Scan(
			&s.ID, &s.RepoName, &s.RepoOwner, &s.RepoFullName,
			&s.StarsCount, &s.ForksCount, &s.ContributorsCount,
			&statsDate, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		s.StatsDate = statsDate
		s.CreatedAt = createdAt
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
