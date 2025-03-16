package entity

import "time"

type RepoDataPoint struct {
	Time         time.Time `json:"time"`
	RepoURL      string    `json:"repo_url"`
	Stars        int       `json:"stargazers_count"`
	Forks        int       `json:"forks_count"`
	Contributors int       `json:"contributors"`
}
