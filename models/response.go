package models

type Usage struct {
	PromptTokens int
	CompletionTokens int
	TotalTokens int
}

type GenerateResponse struct {
	Message Message
	Usage Usage
}