package types

// OpenAIMessage represents a message in the OpenAI chat completion
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

// OpenAIChoice represents a single choice in OpenAI response
type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}
