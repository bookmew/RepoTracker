package util

import (
	"log"
	"time"

	"RepoTracker/src/repository"
)

type RepoStatsScheduler struct {
	repoStatsRepo *repository.RepoStatsRepository
	tokenRepo     *repository.TokenRepository
	interval      time.Duration
	stopChan      chan struct{}
}

func NewRepoStatsScheduler(
	repoStatsRepo *repository.RepoStatsRepository,
	tokenRepo *repository.TokenRepository,
	interval time.Duration,
) *RepoStatsScheduler {
	return &RepoStatsScheduler{
		repoStatsRepo: repoStatsRepo,
		tokenRepo:     tokenRepo,
		interval:      interval,
		stopChan:      make(chan struct{}),
	}
}

func (s *RepoStatsScheduler) Start() {
	ticker := time.NewTicker(s.interval)
	go func() {
		s.fetchAllRepoData()
		
		for {
			select {
			case <-ticker.C:
				s.fetchAllRepoData()
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
	
	log.Printf("Repository data scheduler started with interval: %v", s.interval)
}

func (s *RepoStatsScheduler) Stop() {
	close(s.stopChan)
	log.Println("Repository data scheduler stopped")
}

func (s *RepoStatsScheduler) fetchAllRepoData() {
	tokens, err := s.tokenRepo.GetAll()
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		return
	}

	uniqueRepos := make(map[string]bool)

	for _, token := range tokens {
		if token.RepoURL == "" {
			continue
		}
		
		uniqueRepos[token.RepoURL] = true
	}
	
	for repoURL := range uniqueRepos {
		log.Printf("Fetching data for repository: %s", repoURL)
		
		dataPoint, err := FetchRepoData(repoURL)
		if err != nil {
			log.Printf("Error fetching data for %s: %v", repoURL, err)
			continue
		}
		
		err = s.repoStatsRepo.SaveStats(repoURL, dataPoint)
		if err != nil {
			log.Printf("Error saving data for %s: %v", repoURL, err)
			continue
		}
		
		log.Printf("Successfully updated data for %s: Stars=%d, Forks=%d, Contributors=%d",
			repoURL, dataPoint.Stars, dataPoint.Forks, dataPoint.Contributors)
	}
} 