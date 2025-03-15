package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Repository 结构体用于存储仓库基本信息
type Repository struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	URL         string    `json:"html_url"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Watchers    int       `json:"watchers_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Commit 结构体用于存储提交信息
type Commit struct {
	SHA     string `json:"sha"`
	Commit  struct {
		Message string `json:"message"`
		Author  struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

// Issue 结构体用于存储问题信息
type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
}

func main() {
	// 设置GitHub仓库信息
	owner := "web3-ai-agent" // 仓库所有者
	repo := "aiagent"        // 仓库名称

	// 获取仓库基本信息
	repository, err := getRepository(owner, repo)
	if err != nil {
		log.Fatalf("获取仓库信息失败: %v", err)
	}

	fmt.Printf("仓库名称: %s\n", repository.FullName)
	fmt.Printf("描述: %s\n", repository.Description)
	fmt.Printf("URL: %s\n", repository.URL)
	fmt.Printf("星标数: %d\n", repository.Stars)
	fmt.Printf("分支数: %d\n", repository.Forks)
	fmt.Printf("观察者数: %d\n", repository.Watchers)
	fmt.Printf("创建时间: %s\n", repository.CreatedAt.Format("2006-01-02"))
	fmt.Printf("更新时间: %s\n", repository.UpdatedAt.Format("2006-01-02"))

	// 获取最近的提交
	commits, err := getCommits(owner, repo, 5)
	if err != nil {
		log.Printf("获取提交信息失败: %v", err)
	} else {
		fmt.Println("\n最近的提交:")
		for i, commit := range commits {
			fmt.Printf("%d. [%s] %s by %s (%s)\n", i+1, commit.SHA[:7], commit.Commit.Message, commit.Commit.Author.Name, commit.Commit.Author.Date.Format("2006-01-02"))
		}
	}

	// 获取最近的问题
	issues, err := getIssues(owner, repo, 5)
	if err != nil {
		log.Printf("获取问题信息失败: %v", err)
	} else {
		fmt.Println("\n最近的问题:")
		for i, issue := range issues {
			fmt.Printf("%d. #%d %s [%s] by %s (%s)\n", i+1, issue.Number, issue.Title, issue.State, issue.User.Login, issue.CreatedAt.Format("2006-01-02"))
		}
	}
}

// getRepository 获取仓库基本信息
func getRepository(owner, repo string) (*Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置GitHub API所需的头部
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	// 如果有GitHub令牌，添加到请求中
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}
	
	// 解析响应
	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}
	
	return &repository, nil
}

// getCommits 获取仓库的提交历史
func getCommits(owner, repo string, limit int) ([]Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=%d", owner, repo, limit)
	
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置GitHub API所需的头部
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	// 如果有GitHub令牌，添加到请求中
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}
	
	// 解析响应
	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}
	
	return commits, nil
}

// getIssues 获取仓库的问题列表
func getIssues(owner, repo string, limit int) ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?per_page=%d&state=all", owner, repo, limit)
	
	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置GitHub API所需的头部
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	// 如果有GitHub令牌，添加到请求中
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}
	
	// 解析响应
	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}
	
	return issues, nil
} 