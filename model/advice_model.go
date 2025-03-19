package model

import (
	"fmt"
	"regexp"
)

type ReadmeAdviceRequest struct {
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Content string `json:"content,omitempty"`
}

type ReadmeAdviceResponse struct {
	Evaluation EvaluationResponse `json:"evaluation"`
	Advice     string             `json:"advice"`
	Content    string             `json:"content,omitempty"`
	NewReadme  string             `json:"new_readme"`
	Error      string             `json:"error,omitempty"`
}

func GenerateReadmeAdvice(content string, scores EvaluationResponse) (string, error) {
	advicePrompt := `以下のREADMEに対する評価結果に基づいて、改善のためのアドバイスと改善されたREADMEの例を提供してください。
評価結果:
明確さ: ` + fmt.Sprintf("%d", scores.Clarity) + `/10
完全性: ` + fmt.Sprintf("%d", scores.Completeness) + `/10
構造化: ` + fmt.Sprintf("%d", scores.Structure) + `/10
例示: ` + fmt.Sprintf("%d", scores.Examples) + `/10
可読性: ` + fmt.Sprintf("%d", scores.Readability) + `/10
合計: ` + fmt.Sprintf("%d", scores.TotalScore) + `/50

元のREADME:
"""
` + content + `
"""

以下の形式で回答してください:

## 改善アドバイス
(各評価カテゴリに対する具体的なアドバイスを箇条書きで記載)
(各評価カテゴリに対する点数は含めないでください)

## 改善されたREADMEの例
(元のREADMEを改善した例を記載)

「改善アドバイス」と「改善されたREADMEの例」の2つのセクションを必ず含めてください。`

	adviceText, err := GenerateGeminiResponse(advicePrompt)
	if err != nil {
		return "", err
	}

	return adviceText, nil
}

func GenerateReadmeAdviceWithImprovement(content string) (ReadmeAdviceResponse, error) {
	scores, err := EvaluateReadme(content)
	if err != nil {
		return ReadmeAdviceResponse{}, err
	}

	advice, err := GenerateReadmeAdvice(content, scores)
	if err != nil {
		return ReadmeAdviceResponse{}, err
	}

	improvedReadme := extractImprovedReadme(advice)

	return ReadmeAdviceResponse{
		Evaluation: scores,
		Advice:     advice,
		Content:    content,
		NewReadme:  improvedReadme,
	}, nil
}

func extractImprovedReadme(advice string) string {
	pattern := regexp.MustCompile(`(?s)## 改善されたREADMEの例\s*(.+)$`)
	matches := pattern.FindStringSubmatch(advice)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
