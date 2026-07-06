package gemini

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/YashShekhawat/fusion/models"
)

const defaultBaseURL = "https://generativelanguage.googleapis.com"

type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type GeminiDriver struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func prepareConfig(config *Config) error {
	if strings.TrimSpace(config.APIKey) == "" {
		return fmt.Errorf("gemini: api key is required")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	return nil
}

func New(config Config) (*GeminiDriver, error) {

	if err := prepareConfig(&config); err != nil {
		return nil, err
	}

	return &GeminiDriver{
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}, nil
}

func (d *GeminiDriver) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {

	const generateContentEndpoint = "/v1beta/models/%s:generateContent"
	endpoint := fmt.Sprintf(generateContentEndpoint, req.Model)

	// Convert to Gemini request format
	geminiReq := buildGeminiRequest(req)

	// Send request to gemini
	var geminiResp generateResponse
	if err := d.sendRequest(ctx, endpoint, geminiReq, &geminiResp); err != nil {
		return models.GenerateResponse{}, err
	}

	// Convert gemini response into Fusion response
	fusionResp, err := buildFusionResponse(geminiResp)
	if err != nil {
		return models.GenerateResponse{}, fmt.Errorf("gemini: build fusion response: %w", err)
	}

	return fusionResp, nil
}

func (d *GeminiDriver) Name() string {
	return "gemini"
}
