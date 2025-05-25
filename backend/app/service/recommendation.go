package service

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type (
	// Recommendation トレーニングメニュー提案のサービスインターフェース
	Recommendation interface {
		ProposeTrainingMenu(goal string, parts []string, experience string, time int) (string, error)
	}

	// RecommendationImpl トレーニングメニュー提案のサービス実装
	RecommendationImpl struct {
		openAIClient *openai.Client
	}
)

// コンストラクタ: 環境変数等からAPIキーを読み込んでクライアントを初期化
func NewRecommendation() Recommendation {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return &RecommendationImpl{
		openAIClient: client,
	}
}

// トレーニングメニュー提案ロジック
func (s *RecommendationImpl) ProposeTrainingMenu(goal string, parts []string, experience string, time int) (string, error) {
	prompt := buildPrompt(goal, parts, experience, time)

	resp, err := s.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-3.5-turbo",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "あなたはプロのパーソナルトレーナーです。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   800, // 必要に応じて調整
			Temperature: 0.7, // ランダム性
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	// モデルからの提案結果(文章)を取得
	result := resp.Choices[0].Message.Content
	return result, nil
}

// プロンプトを組み立てる関数
func buildPrompt(goal string, parts []string, experience string, time int) string {
	return fmt.Sprintf(
		`
		トレーニング目的: %s
		対象部位: %v
		トレーニング経験: %s
		確保できる時間: %d分

		上記の条件に合わせて、適切な筋トレメニューを提案してください。具体的な種目、回数、セット数、インターバルなども含めて提示をお願いします。`,
		goal, parts, experience, time,
	)
}
