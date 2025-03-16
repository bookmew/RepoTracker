package service

import (
	"RepoTracker/src/entity"
	"RepoTracker/src/repository"
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
