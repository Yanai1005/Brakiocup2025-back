package model

import (
	"regexp"
	"strconv"
)

type EvaluationRequest struct {
	Content string `json:"content"`
}

type EvaluationResponse struct {
	Clarity      int    `json:"clarity"`
	Completeness int    `json:"completeness"`
	Structure    int    `json:"structure"`
	Examples     int    `json:"examples"`
	Readability  int    `json:"readability"`
	TotalScore   int    `json:"total_score"`
	Error        string `json:"error,omitempty"`
}

func EvaluateReadme(content string) (EvaluationResponse, error) {
	// 評価プロンプトを構築
	evaluationPrompt := `以下のREADMEファイルを評価してください。評価基準は次のとおりです：

1. 明確さ: プロジェクトの目的を明確に説明されているか（0〜10点）
2. 完全性:  環境構築、インストール手順、使用方法などの必要な情報がすべて含まれているか（0〜10点）
3. 構造化: 情報が論理的に整理されているか（0〜10点）
4. 例示: 使用例やコード例が含まれているか（0〜10点）
5. 可読性: 文章が読みやすいか（0〜10点）

回答は以下の形式で返してください：
明確さ: X点
完全性: X点
構造化: X点
例示: X点
可読性: X点
合計: X点

評価対象のテキスト:
"""
` + content + `
"""

コメントや追加の説明は不要です。点数のみを返してください。`

	evaluationText, err := GenerateGeminiResponse(evaluationPrompt)
	if err != nil {
		return EvaluationResponse{}, err
	}

	return extractScores(evaluationText), nil
}

func extractScores(text string) EvaluationResponse {
	var response EvaluationResponse

	clarityPattern := regexp.MustCompile(`明確さ:\s*(\d+)`)
	completenessPattern := regexp.MustCompile(`完全性:\s*(\d+)`)
	structurePattern := regexp.MustCompile(`構造化:\s*(\d+)`)
	examplesPattern := regexp.MustCompile(`例示:\s*(\d+)`)
	readabilityPattern := regexp.MustCompile(`可読性:\s*(\d+)`)
	totalPattern := regexp.MustCompile(`合計(?:点数)?:\s*(\d+)`)

	if matches := clarityPattern.FindStringSubmatch(text); len(matches) > 1 {
		response.Clarity, _ = strconv.Atoi(matches[1])
	}

	if matches := completenessPattern.FindStringSubmatch(text); len(matches) > 1 {
		response.Completeness, _ = strconv.Atoi(matches[1])
	}

	if matches := structurePattern.FindStringSubmatch(text); len(matches) > 1 {
		response.Structure, _ = strconv.Atoi(matches[1])
	}

	if matches := examplesPattern.FindStringSubmatch(text); len(matches) > 1 {
		response.Examples, _ = strconv.Atoi(matches[1])
	}

	if matches := readabilityPattern.FindStringSubmatch(text); len(matches) > 1 {
		response.Readability, _ = strconv.Atoi(matches[1])
	}

	if matches := totalPattern.FindStringSubmatch(text); len(matches) > 1 {
		response.TotalScore, _ = strconv.Atoi(matches[1])
	} else {
		response.TotalScore = response.Clarity + response.Completeness +
			response.Structure + response.Examples + response.Readability
	}

	return response
}
