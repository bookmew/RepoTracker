package repository

import (
	"database/sql"
	"time"

	"RepoTracker/src/entity"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) GetAll() ([]entity.Token, error) {
	query := `
		SELECT id, symbol, name, address, metadata, repo_url, created_at, updated_at
		FROM tokens
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []entity.Token
	for rows.Next() {
		var t entity.Token
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&t.ID, &t.Symbol, &t.Name, &t.Address,
			&t.Metadata, &t.RepoURL, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		t.CreatedAt = createdAt
		t.UpdatedAt = updatedAt
		tokens = append(tokens, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}
