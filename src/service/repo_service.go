package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github-repo-tracker/src/model"
	"github-repo-tracker/src/repository"
	"github-repo-tracker/src/util"
)

type RepoService struct {
	repo        repository.RepoRepository
	githubClient *util.GitHubClient
}

func NewRepoService(repo repository.RepoRepository) *RepoService {
	return &RepoService{
		repo:        repo,
		githubClient: util.NewGitHubClient(),
	}
}

func (s *RepoService) GetAllRepositories(ctx context.Context) ([]model.Repository, error) {
	return s.repo.GetAll(ctx)
}

func (s *RepoService) GetRepositoryByID(ctx context.Context, id int) (*model.Repository, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RepoService) GetRepositoryByFullName(ctx context.Context, fullName string) (*model.Repository, error) {
	return s.repo.GetByFullName(ctx, fullName)
}

func (s *RepoService) CreateRepository(ctx context.Context, repo *model.Repository) error {
	if repo.FullName == "" && repo.Owner != "" && repo.Name != "" {
		repo.FullName = repo.Owner + "/" + repo.Name
	}

	if repo.CreatedAt.IsZero() {
		repo.CreatedAt = time.Now()
	}

	repo.LastUpdated = time.Now()
	
	return s.repo.Create(ctx, repo)
}

func (s *RepoService) UpdateRepositoryStats(ctx context.Context, id int, stats *model.RepositoryStats) error {
	return s.repo.UpdateStats(ctx, id, stats)
}

func (s *RepoService) FetchAndUpdateRepositoryStats(ctx context.Context, fullName string) error {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return nil
	}
	
	owner := parts[0]
	name := parts[1]
	
	stats, err := s.githubClient.GetRepositoryStats(owner, name)
	if err != nil {
		return err
	}
	
	repo, err := s.repo.GetByFullName(ctx, fullName)
	if err != nil {
		newRepo := &model.Repository{
			Name:             name,
			Owner:            owner,
			FullName:         fullName,
			StarsCount:       stats.StarsCount,
			ForksCount:       stats.ForksCount,
			ContributorsCount: stats.ContributorsCount,
			LastUpdated:      time.Now(),
			CreatedAt:        time.Now(),
		}
		
		return s.repo.Create(ctx, newRepo)
	}

	repo.StarsCount = stats.StarsCount
	repo.ForksCount = stats.ForksCount
	repo.ContributorsCount = stats.ContributorsCount
	repo.LastUpdated = time.Now()
	
	return s.repo.Update(ctx, repo)
}

func (s *RepoService) StartPeriodicStatsFetcher(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.updateAllRepositoriesStats(ctx)
			}
		}
	}()
}

func (s *RepoService) updateAllRepositoriesStats(ctx context.Context) {
	repos, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		return
	}
	
	for _, repo := range repos {
		err := s.FetchAndUpdateRepositoryStats(ctx, repo.FullName)
		if err != nil {
			log.Printf("Error updating stats for %s: %v", repo.FullName, err)
		} else {
			log.Printf("Updated stats for %s", repo.FullName)
		}

		time.Sleep(1 * time.Second)
	}
} 