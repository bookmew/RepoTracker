package repository

import (
	"context"
	"time"

	"github-repo-tracker/src/model"

	"github.com/jackc/pgx/v4"
)

type RepoRepository interface {
	GetAll(ctx context.Context) ([]model.Repository, error)

	GetByID(ctx context.Context, id int) (*model.Repository, error)

	GetByFullName(ctx context.Context, fullName string) (*model.Repository, error)

	Create(ctx context.Context, repo *model.Repository) error

	Update(ctx context.Context, repo *model.Repository) error

	UpdateStats(ctx context.Context, id int, stats *model.RepositoryStats) error
}

type PostgresRepoRepository struct {
	db *pgx.Conn
}

func NewPostgresRepoRepository(db *pgx.Conn) *PostgresRepoRepository {
	return &PostgresRepoRepository{db: db}
}

func (r *PostgresRepoRepository) GetAll(ctx context.Context) ([]model.Repository, error) {
	query := `SELECT id, name, owner, full_name, stars_count, forks_count, contributors_count, last_updated, created_at 
			  FROM repositories`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var repos []model.Repository
	for rows.Next() {
		var repo model.Repository
		err := rows.Scan(
			&repo.ID,
			&repo.Name,
			&repo.Owner,
			&repo.FullName,
			&repo.StarsCount,
			&repo.ForksCount,
			&repo.ContributorsCount,
			&repo.LastUpdated,
			&repo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	
	return repos, nil
}

func (r *PostgresRepoRepository) GetByID(ctx context.Context, id int) (*model.Repository, error) {
	query := `SELECT id, name, owner, full_name, stars_count, forks_count, contributors_count, last_updated, created_at 
			  FROM repositories WHERE id = $1`
	
	var repo model.Repository
	err := r.db.QueryRow(ctx, query, id).Scan(
		&repo.ID,
		&repo.Name,
		&repo.Owner,
		&repo.FullName,
		&repo.StarsCount,
		&repo.ForksCount,
		&repo.ContributorsCount,
		&repo.LastUpdated,
		&repo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &repo, nil
}

func (r *PostgresRepoRepository) GetByFullName(ctx context.Context, fullName string) (*model.Repository, error) {
	query := `SELECT id, name, owner, full_name, stars_count, forks_count, contributors_count, last_updated, created_at 
			  FROM repositories WHERE full_name = $1`
	
	var repo model.Repository
	err := r.db.QueryRow(ctx, query, fullName).Scan(
		&repo.ID,
		&repo.Name,
		&repo.Owner,
		&repo.FullName,
		&repo.StarsCount,
		&repo.ForksCount,
		&repo.ContributorsCount,
		&repo.LastUpdated,
		&repo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &repo, nil
}

func (r *PostgresRepoRepository) Create(ctx context.Context, repo *model.Repository) error {
	query := `INSERT INTO repositories (name, owner, full_name, stars_count, forks_count, contributors_count, last_updated, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  RETURNING id`
	
	repo.CreatedAt = time.Now()
	repo.LastUpdated = time.Now()
	
	return r.db.QueryRow(ctx, query,
		repo.Name,
		repo.Owner,
		repo.FullName,
		repo.StarsCount,
		repo.ForksCount,
		repo.ContributorsCount,
		repo.LastUpdated,
		repo.CreatedAt,
	).Scan(&repo.ID)
}

func (r *PostgresRepoRepository) Update(ctx context.Context, repo *model.Repository) error {
	query := `UPDATE repositories
			  SET name = $1, owner = $2, full_name = $3, stars_count = $4, forks_count = $5, contributors_count = $6, last_updated = $7
			  WHERE id = $8`
	
	repo.LastUpdated = time.Now()
	
	_, err := r.db.Exec(ctx, query,
		repo.Name,
		repo.Owner,
		repo.FullName,
		repo.StarsCount,
		repo.ForksCount,
		repo.ContributorsCount,
		repo.LastUpdated,
		repo.ID,
	)
	
	return err
}

func (r *PostgresRepoRepository) UpdateStats(ctx context.Context, id int, stats *model.RepositoryStats) error {
	query := `UPDATE repositories
			  SET stars_count = $1, forks_count = $2, contributors_count = $3, last_updated = $4
			  WHERE id = $5`
	
	_, err := r.db.Exec(ctx, query,
		stats.StarsCount,
		stats.ForksCount,
		stats.ContributorsCount,
		time.Now(),
		id,
	)
	
	return err
} 