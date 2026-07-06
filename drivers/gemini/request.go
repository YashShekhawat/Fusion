package gemini

import (
	"strings"

	"github.com/YashShekhawat/fusion/models"
)

type generateRequest struct {
	SystemInstruction *content  `json:"systemInstruction,omitempty"`
	Contents          []content `json:"contents"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

func buildGeminiRequest(req models.GenerateRequest) generateRequest {
	geminiReq := generateRequest{
		Contents: make([]content, 0, len(req.Messages)),
	}

	var systemInstructions []string

	for _, msg := range req.Messages {

		if strings.TrimSpace(msg.Content) == "" {
			continue
		}

		text := strings.TrimSpace(msg.Content)
		if text == "" {
			continue
		}

		switch msg.Role {

		case models.RoleSystem:
			systemInstructions = append(systemInstructions, text)

		default:
			geminiReq.Contents = append(geminiReq.Contents, content{
				Role: string(msg.Role),
				Parts: []part{
					{
						Text: text,
					},
				},
			})
		}
	}

	// Merge all system messages into a single Gemini systemInstruction.
	if len(systemInstructions) > 0 {
		geminiReq.SystemInstruction = &content{
			Parts: []part{
				{
					Text: strings.Join(systemInstructions, "\n\n"),
				},
			},
		}
	}

	return geminiReq
}
