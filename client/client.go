package client

import (
	"context"
	"fmt"

	"github.com/YashShekhawat/fusion/drivers"
	"github.com/YashShekhawat/fusion/drivers/gemini"
	"github.com/YashShekhawat/fusion/drivers/openai"
	"github.com/YashShekhawat/fusion/middleware"
	"github.com/YashShekhawat/fusion/models"
	"github.com/YashShekhawat/fusion/registry"
)

type Client struct {
	registry    *registry.Registry
	middlewares []middleware.Middleware
	provider    Provider
	apiKey      string
}

func New(opts ...Option) (*Client, error) {
	client := &Client{
		registry: registry.New(),
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if err := client.registerProvider(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Register(d drivers.Driver) error {
	return c.registry.Register(d)
}

func (c *Client) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {

	driver, err := c.registry.Get(string(c.provider))
	if err != nil {
		return models.GenerateResponse{}, fmt.Errorf("failed to get driver: %w", err)
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		driver = c.middlewares[i](driver)
	}

	return driver.Generate(ctx, req)
}

func (c *Client) registerProvider() error {
	switch c.provider {

	case ProviderGemini:
		driver, err := gemini.New(gemini.Config{
			APIKey: c.apiKey,
		})
		if err != nil {
			return err
		}
		return c.Register(driver)

	case ProviderOpenAI:
		driver, err := openai.New(openai.Config{
			APIKey: c.apiKey,
		})
		if err != nil {
			return err
		}
		return c.Register(driver)

	default:
		return fmt.Errorf("unsupported provider %q", c.provider)
	}
}
