package client

import (
	"context"
	"fmt"

	"github.com/YashShekhawat/fusion/drivers"
	"github.com/YashShekhawat/fusion/models"
	"github.com/YashShekhawat/fusion/registry"
)

type Client struct {
	registry *registry.Registry
}

func New() *Client {
	return &Client{registry: registry.New()}
}

func (c *Client) Register(d drivers.Drivers) error {
	return c.registry.Register(d)
}

func (c *Client) Generate(ctx context.Context, driverName string, req models.GenerateRequest) (models.GenerateResponse, error) {
	fmt.Println("Client: Looking up driver")
	d, err := c.registry.Get(driverName)
	if err != nil {
		return models.GenerateResponse{}, fmt.Errorf("failed to get driver: %w", err)
	}
	fmt.Println("Client: Calling driver.Generate()")
	return d.Generate(ctx, req)
}
