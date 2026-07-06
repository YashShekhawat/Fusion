package client

import "github.com/YashShekhawat/fusion/middleware"

type Option func(*Client) error

func WithMiddleware(mw ...middleware.Middleware) Option {
	return func(c *Client) error {
		c.middlewares = append(c.middlewares, mw...)
		return nil
	}
}

func WithProvider(provider Provider) Option {
	return func(c *Client) error {
		c.provider = provider
		return nil
	}
}

func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}
