package openai

import "github.com/YashShekhawat/fusion/models"

type generateRequest struct {
	Model string         `json:"model"`
	Input []inputMessage `json:"input"`
}

type inputMessage struct {
	Role    string         `json:"role"`
	Content []inputContent `json:"content"`
}

type inputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func buildOpenAIRequest(req models.GenerateRequest) generateRequest {
	openAIReq := generateRequest{
		Model: req.Model,
		Input: make([]inputMessage, 0, len(req.Messages)),
	}

	for _, msg := range req.Messages {
		openAIReq.Input = append(openAIReq.Input, inputMessage{
			Role: string(msg.Role),
			Content: []inputContent{
				{
					Type: "input_text",
					Text: msg.Content,
				},
			},
		})
	}
	return openAIReq
}
