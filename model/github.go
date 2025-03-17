package model

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type GitHubRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type GitHubResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func FetchReadme(owner, repo string) (string, error) {
	if owner == "" || repo == "" {
		return "", fmt.Errorf("owner と repo は必須です")
	}

	client := github.NewClient(nil)

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
