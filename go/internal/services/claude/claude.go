package claude

import (
	"context"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type ClaudeClient struct {
	client *anthropic.Client
}

func NewClaudeClient() *ClaudeClient {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return nil
	}
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &ClaudeClient{client: &client}
}

func (c *ClaudeClient) FindAnimeByDescription(prompt string) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}
	params := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(
				`You are an anime expert. Given a description, recommend 3 anime titles that match. Respond ONLY with a valid JSON array in this format:
[
  {"title":"Anime Name","reason":"Why it matches","confidence":0.9}
]
` + prompt)),
		},
	}
	message, err := c.client.Messages.New(context.TODO(), params)
	if err != nil {
		return "", err
	}
	if len(message.Content) > 0 {
		return message.Content[0].Text, nil
	}
	return "", nil
}
