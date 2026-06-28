package mock

import (
	"context"
	"fmt"

	"github.com/YashShekhawat/fusion/models"
)

type Mock struct{}

func New() *Mock {
	return &Mock{}
}

func (m *Mock) Name() string {
	return "mock"
}

func (m *Mock) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {
	fmt.Println("MockDriver: Generating response")
	return models.GenerateResponse{
		Message: models.Message{
			Role:    models.RoleAssistant,
			Content: "This is a Mock response",
		},
		Usage: models.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}
