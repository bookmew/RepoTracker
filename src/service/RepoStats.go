package service

import (
	"RepoTracker/src/entity"
	"RepoTracker/src/repository"
	"RepoTracker/src/util"
)

type RepoStatsService struct {
	repo *repository.RepoStatsRepository
}

func NewRepoStatsService(repo *repository.RepoStatsRepository) *RepoStatsService {
	return &RepoStatsService{repo: repo}
}

func (s *RepoStatsService) GetByMint(repoFullName string) ([]entity.RepoStats, error) {
	return s.repo.GetByMint(repoFullName)
}

func (s *RepoStatsService) GetAll() ([]entity.RepoStats, error) {
	return s.repo.GetAll()
}

func (s *RepoStatsService) SaveStats(repoURL string) error {
	// Fetch repository data from GitHub API
	dataPoint, err := util.FetchRepoData(repoURL)
	if err != nil {
		return err
	}

	// Save data to database
	return s.repo.SaveStats(repoURL, dataPoint)
}
