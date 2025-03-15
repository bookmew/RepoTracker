package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github-repo-tracker/src/model"
)

type GitHubClient struct {
	client    *http.Client
	token     string
	baseURL   string
	rateLimit int
}

func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	return &GitHubClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		token:     token,
		baseURL:   "https://api.github.com",
		rateLimit: 5000,
	}
}

func (c *GitHubClient) GetRepositoryStats(owner, repo string) (*model.RepositoryStats, error) {
	repoURL := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)
	req, err := http.NewRequest("GET", repoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var repoData struct {
		StargazersCount int `json:"stargazers_count"`
		ForksCount      int `json:"forks_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	contributorsURL := fmt.Sprintf("%s/repos/%s/%s/contributors?per_page=1&anon=true", c.baseURL, owner, repo)
	req, err = http.NewRequest("GET", contributorsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating contributors request: %v", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making contributors request: %v", err)
	}
	defer resp.Body.Close()

	contributorsCount := 0
	if resp.StatusCode == http.StatusOK {
		if link := resp.Header.Get("Link"); link != "" {
			if lastLink := extractLastLink(link); lastLink != "" {
				if page := extractPageNumber(lastLink); page > 0 {
					contributorsCount = page
				}
			}
		}
	}

	return &model.RepositoryStats{
		StarsCount:       repoData.StargazersCount,
		ForksCount:       repoData.ForksCount,
		ContributorsCount: contributorsCount,
	}, nil
}

func extractLastLink(link string) string {
	for _, part := range []string{link} {
		if part != "" && part[len(part)-6:] == "last\"" {
			for i := 0; i < len(part); i++ {
				if part[i] == '<' {
					for j := i + 1; j < len(part); j++ {
						if part[j] == '>' {
							return part[i+1 : j]
						}
					}
				}
			}
		}
	}
	return ""
}

func extractPageNumber(url string) int {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '=' {
			pageStr := url[i+1:]
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				return 0
			}
			return page
		}
	}
	return 0
} 