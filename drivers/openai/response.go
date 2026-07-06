package openai

import (
	"fmt"

	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
	"github.com/YashShekhawat/fusion/models"
)

type generateResponse struct {
	Output []outputMessage `json:"output"`
	Usage  usage           `json:"usage"`
}

type outputMessage struct {
	Role    string          `json:"role"`
	Content []outputContent `json:"content"`
}

type outputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func buildFusionResponse(resp generateResponse) (models.GenerateResponse, error) {

	//validate output
	if len(resp.Output) == 0 {
		return models.GenerateResponse{}, fmt.Errorf("openai: invalid response payload: %w", fusionerrors.ErrInvalidResponse)
	}
	// validate content
	if len(resp.Output[0].Content) == 0 {
		return models.GenerateResponse{}, fmt.Errorf("openai: invalid response payload: %w", fusionerrors.ErrInvalidResponse)
	}

	fusionResp := models.GenerateResponse{
		Message: models.Message{
			Role:    models.Role(resp.Output[0].Role),
			Content: resp.Output[0].Content[0].Text,
		},
		Usage: models.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	return fusionResp, nil
}
