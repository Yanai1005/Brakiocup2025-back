package model

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHubRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type GitHubResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// ユーザー分析リクエスト
type UserAnalysisRequest struct {
	Username string `json:"username"`
	MaxRepos int    `json:"max_repos"` // 最大分析リポジトリ数
}

// ユーザー分析レスポンス
type UserAnalysisResponse struct {
	Username        string               `json:"username"`
	RepositoryCount int                  `json:"repository_count"`
	AnalyzedCount   int                  `json:"analyzed_count"`
	AverageScores   EvaluationResponse   `json:"average_scores"`
	RepoAnalyses    []RepositoryAnalysis `json:"repo_analyses"`
	Error           string               `json:"error,omitempty"`
}

type RepositoryAnalysis struct {
	RepoName    string             `json:"repo_name"`
	Scores      EvaluationResponse `json:"scores"`
	HasReadme   bool               `json:"has_readme"`
	Description string             `json:"description,omitempty"`
}

func FetchReadme(owner, repo string) (string, error) {
	if owner == "" || repo == "" {
		return "", fmt.Errorf("owner と repo は必須です")
	}

	var client *github.Client
	token := os.Getenv("BACK_PAT")

	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	ctx := context.Background()

	readme, _, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if err != nil {
		return "", fmt.Errorf("README の取得に失敗: %v", err)
	}

	content, err := readme.GetContent()
	if err != nil {
		return "", fmt.Errorf("README コンテンツの抽出に失敗: %v", err)
	}

	return content, nil
}

func FetchUserRepositories(username string, maxRepos int) ([]*github.Repository, error) {
	if username == "" {
		return nil, fmt.Errorf("username は必須です")
	}

	if maxRepos <= 0 || maxRepos > 10 {
		maxRepos = 10
	}

	var client *github.Client
	token := os.Getenv("BACK_PAT")

	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	ctx := context.Background()

	opt := &github.RepositoryListOptions{
		Type:        "public",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: maxRepos},
	}

	repos, _, err := client.Repositories.List(ctx, username, opt)
	if err != nil {
		return nil, fmt.Errorf("リポジトリの取得に失敗: %v", err)
	}

	return repos, nil
}

func AnalyzeUserRepositories(username string, maxRepos int) (UserAnalysisResponse, error) {
	response := UserAnalysisResponse{
		Username:      username,
		RepoAnalyses:  []RepositoryAnalysis{},
		AverageScores: EvaluationResponse{},
	}

	repos, err := FetchUserRepositories(username, maxRepos)
	if err != nil {
		return response, err
	}

	response.RepositoryCount = len(repos)

	totalScores := EvaluationResponse{}
	analyzedCount := 0

	for _, repo := range repos {
		repoAnalysis := RepositoryAnalysis{
			RepoName:    *repo.Name,
			HasReadme:   false,
			Description: getRepoDescription(repo),
		}
		readmeContent, err := FetchReadme(username, *repo.Name)
		if err != nil {
			response.RepoAnalyses = append(response.RepoAnalyses, repoAnalysis)
			continue
		}

		repoAnalysis.HasReadme = true
		scores, err := EvaluateReadme(readmeContent)
		if err != nil {
			response.RepoAnalyses = append(response.RepoAnalyses, repoAnalysis)
			continue
		}

		repoAnalysis.Scores = scores
		response.RepoAnalyses = append(response.RepoAnalyses, repoAnalysis)

		totalScores.Clarity += scores.Clarity
		totalScores.Completeness += scores.Completeness
		totalScores.Structure += scores.Structure
		totalScores.Examples += scores.Examples
		totalScores.Readability += scores.Readability
		totalScores.TotalScore += scores.TotalScore
		analyzedCount++
	}

	response.AnalyzedCount = analyzedCount
	if analyzedCount > 0 {
		response.AverageScores = EvaluationResponse{
			Clarity:      totalScores.Clarity / analyzedCount,
			Completeness: totalScores.Completeness / analyzedCount,
			Structure:    totalScores.Structure / analyzedCount,
			Examples:     totalScores.Examples / analyzedCount,
			Readability:  totalScores.Readability / analyzedCount,
			TotalScore:   totalScores.TotalScore / analyzedCount,
		}
	}

	return response, nil
}

func getRepoDescription(repo *github.Repository) string {
	if repo.Description != nil {
		return *repo.Description
	}
	return ""
}
