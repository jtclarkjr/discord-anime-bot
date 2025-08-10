package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"discord-anime-bot/internal/types"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

// FindAnimeByDescription uses OpenAI to find anime recommendations based on description
func FindAnimeByDescription(description, apiKey string) ([]types.OpenAIRecommendation, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI is not configured. Please set OPENAI_API_KEY environment variable")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	prompt := fmt.Sprintf(`Based on this description: "%s"

Please recommend anime titles that match this description. Return your response as a JSON array of objects with the following structure:
[
  {
    "title": "Exact anime title",
    "reason": "Brief explanation of why this matches",
    "confidence": 0.95
  }
]

Guidelines:
- Return 1-3 recommendations
- Use exact anime titles (romaji or English)
- Confidence should be between 0.0 and 1.0
- Only return valid JSON, no other text
- Focus on popular/well-known anime`, description)

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT5, // Uses the latest model
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	var recommendations []types.OpenAIRecommendation
	if err := json.Unmarshal([]byte(content), &recommendations); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return recommendations, nil
}
