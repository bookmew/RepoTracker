package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"RepoTracker/src/entity"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}
}

func ExtractOwnerAndRepo(repoURL string) (string, string, error) {
	repoURL = strings.TrimSuffix(repoURL, "/")

	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)
	matches := re.FindStringSubmatch(repoURL)

	if len(matches) < 3 {
		re = regexp.MustCompile(`api\.github\.com/repos/([^/]+)/([^/]+)`)
		matches = re.FindStringSubmatch(repoURL)

		if len(matches) < 3 {
			return "", "", fmt.Errorf("invalid GitHub repository URL: %s", repoURL)
		}
	}

	return matches[1], matches[2], nil
}

func GetAPIURL(repoURL string) (string, error) {
	owner, repo, err := ExtractOwnerAndRepo(repoURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), nil
}

func FetchRepoData(repoURL string) (entity.RepoDataPoint, error) {
	apiURL, err := GetAPIURL(repoURL)
	if err != nil {
		return entity.RepoDataPoint{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return entity.RepoDataPoint{}, err
	}

	authToken := os.Getenv("GITHUB_AUTH_TOKEN")
	if authToken == "" {
		return entity.RepoDataPoint{}, fmt.Errorf("github auth token is not set")
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return entity.RepoDataPoint{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entity.RepoDataPoint{}, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var repoData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return entity.RepoDataPoint{}, err
	}

	stars, _ := repoData["stargazers_count"].(float64)
	forks, _ := repoData["forks_count"].(float64)

	dataPoint := entity.RepoDataPoint{
		Time:    time.Now(),
		RepoURL: repoURL,
		Stars:   int(stars),
		Forks:   int(forks),
	}

	contributorsCount, err := FetchContributorsCount(apiURL)
	if err != nil {
		log.Printf("Warning: Failed to fetch contributors count: %v", err)
		dataPoint.Contributors = 0
	} else {
		dataPoint.Contributors = contributorsCount
	}

	return dataPoint, nil
}

func FetchContributorsCount(apiURL string) (int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contributors?per_page=1", apiURL), nil)
	if err != nil {
		return 0, err
	}

	authToken := os.Getenv("GITHUB_AUTH_TOKEN")
	if authToken == "" {
		return 0, fmt.Errorf("github auth token is not set")
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch contributors: %s", resp.Status)
	}

	linkHeader := resp.Header.Get("Link")

	contributorsCount := 0
	if linkHeader != "" {
		lastLink := ""
		for _, link := range strings.Split(linkHeader, ",") {
			if strings.Contains(link, `rel="last"`) {
				lastLink = link
				break
			}
		}

		if lastLink != "" {
			re := regexp.MustCompile(`page=(\d+)(?:[^&\d]|$)`)
			matches := re.FindStringSubmatch(lastLink)

			if len(matches) > 1 {
				pageNum, err := strconv.Atoi(matches[1])
				if err != nil {
					log.Println("Error converting page number:", err)
				} else {
					contributorsCount = pageNum
				}
			} else {
				log.Println("No page number found in last link")
			}
		}
	} else {
		var contributors []entity.Contributor
		if err := json.NewDecoder(resp.Body).Decode(&contributors); err != nil {
			return 0, err
		}
		contributorsCount = len(contributors)
	}

	return contributorsCount, nil
}
