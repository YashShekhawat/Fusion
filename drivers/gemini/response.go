package gemini

import (
	"fmt"
	"strings"

	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
	"github.com/YashShekhawat/fusion/models"
)

type generateResponse struct {
	Candidates    []candidateMessage `json:"candidates"`
	UsageMetadata usageMetadata      `json:"usageMetadata"`
}

type candidateMessage struct {
	Content messageContent `json:"content"`
}

type messageContent struct {
	Role  string      `json:"role"`
	Parts []partsText `json:"parts"`
}

type partsText struct {
	Text string `json:"text"`
}

type usageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

func buildFusionResponse(resp generateResponse) (models.GenerateResponse, error) {

	//validate output
	if len(resp.Candidates) == 0 {
		return models.GenerateResponse{}, fmt.Errorf("gemini response contained no candidates: %w", fusionerrors.ErrInvalidResponse)
	}
	// validate content
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return models.GenerateResponse{}, fmt.Errorf("gemini response contained no parts: %w", fusionerrors.ErrInvalidResponse)
	}

	var builder strings.Builder
	candidate := resp.Candidates[0]

	for _, part := range candidate.Content.Parts {
		builder.WriteString(part.Text)
	}
	text := strings.TrimSpace(builder.String())

	if text == "" {
		return models.GenerateResponse{}, fmt.Errorf(
			"gemini response contained no content: %w",
			fusionerrors.ErrInvalidResponse,
		)
	}

	role := models.RoleAssistant

	switch candidate.Content.Role {
	case "user":
		role = models.RoleUser
	case "model":
		role = models.RoleAssistant
	}

	fusionResp := models.GenerateResponse{
		Message: models.Message{
			Role:    role,
			Content: text,
		},
		Usage: models.Usage{
			PromptTokens:     resp.UsageMetadata.PromptTokenCount,
			CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      resp.UsageMetadata.TotalTokenCount,
		},
	}

	return fusionResp, nil
}
