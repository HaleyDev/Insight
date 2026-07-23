package agent

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Agent struct {
	llm llms.Model
}

func NewAgent(apiKey, baseURL, model string) (*Agent, error) {
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL(baseURL),
		openai.WithModel(model),
	)
	if err != nil {
		return nil, fmt.Errorf("agent: create llm: %w", err)
	}
	return &Agent{llm: llm}, nil
}

// Chat 执行一次对话。
func (a *Agent) Chat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userMessage),
	}

	resp, err := a.llm.GenerateContent(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("agent: generate content: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("agent: empty response")
	}

	return resp.Choices[0].Content, nil
}
