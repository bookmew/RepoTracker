package entity

import (
	"time"
)

type RepoStats struct {
	ID                int64     `json:"id"`
	RepoName          string    `json:"repo_name"`
	RepoOwner         string    `json:"repo_owner"`
	RepoFullName      string    `json:"repo_full_name"`
	StarsCount        int64     `json:"stars_count"`
	ForksCount        int64     `json:"forks_count"`
	ContributorsCount int64     `json:"contributors_count"`
	StatsDate         time.Time `json:"stats_date"`
	UpdatedAt         time.Time `json:"updated_at"`
}
