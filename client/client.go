package client

import (
	"context"
	"fmt"

	"github.com/YashShekhawat/fusion/drivers"
	"github.com/YashShekhawat/fusion/middleware"
	"github.com/YashShekhawat/fusion/models"
	"github.com/YashShekhawat/fusion/registry"
)

type Client struct {
	registry    *registry.Registry
	middlewares []middleware.Middleware
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

	return client, nil
}

func (c *Client) Register(d drivers.Driver) error {
	return c.registry.Register(d)
}

func (c *Client) Generate(ctx context.Context, driverName string, req models.GenerateRequest) (models.GenerateResponse, error) {
	fmt.Println("Client: Looking up driver")
	d, err := c.registry.Get(driverName)
	if err != nil {
		return models.GenerateResponse{}, fmt.Errorf("failed to get driver: %w", err)
	}
	// Apply configured middleware.
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		d = c.middlewares[i](d)
	}

	fmt.Println("Client: Calling driver.Generate()")
	return d.Generate(ctx, req)
}
