package openai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/YashShekhawat/fusion/models"
)

type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type OpenAIDriver struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

const (
	responsesEndpoint = "/responses"
)

func New(config Config) (*OpenAIDriver, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	return &OpenAIDriver{
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}, nil
}

func (d *OpenAIDriver) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {
	// Convert to OpenAI request format
	openAIReq := buildOpenAIRequest(req)

	// Send request to openai
	var openAIResp generateResponse
	if err := d.sendRequest(ctx, responsesEndpoint, openAIReq, &openAIResp); err != nil {
		return models.GenerateResponse{}, err
	}

	// Convert OpenAI response into Fusion response
	fusionResp, err := buildFusionResponse(openAIResp)
	if err != nil {
		return models.GenerateResponse{}, fmt.Errorf("openai: build fusion response: %w", err)
	}

	return fusionResp, nil
}

func (d *OpenAIDriver) Name() string {
	return "openai"
}
