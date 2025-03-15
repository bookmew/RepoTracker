package model

import (
	"time"
)

// Repository represents a GitHub repository with its statistics
type Repository struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Owner            string    `json:"owner"`
	FullName         string    `json:"full_name"`
	StarsCount       int       `json:"stars_count"`
	ForksCount       int       `json:"forks_count"`
	ContributorsCount int      `json:"contributors_count"`
	LastUpdated      time.Time `json:"last_updated"`
	CreatedAt        time.Time `json:"created_at"`
}

// RepositoryStats represents the statistics of a GitHub repository
type RepositoryStats struct {
	StarsCount       int `json:"stars_count"`
	ForksCount       int `json:"forks_count"`
	ContributorsCount int `json:"contributors_count"`
} 