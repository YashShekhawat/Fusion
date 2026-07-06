package middleware

import (
	"context"
	"time"

	"github.com/YashShekhawat/fusion/drivers"
	"github.com/YashShekhawat/fusion/models"
)

type Logger interface {
	Printf(format string, args ...any)
}

type loggingDriver struct {
	next   drivers.Driver
	logger Logger
}

func Logging(logger Logger) Middleware {
	return func(next drivers.Driver) drivers.Driver {
		return &loggingDriver{
			next:   next,
			logger: logger,
		}
	}
}

func (d *loggingDriver) Name() string {
	return d.next.Name()
}

func (d *loggingDriver) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {
	start := time.Now()

	d.logger.Printf("[%s] Generate started", d.Name())

	resp, err := d.next.Generate(ctx, req)

	duration := time.Since(start)
	if err != nil {
		d.logger.Printf(
			"[%s] Generate failed after %s: %v",
			d.Name(),
			duration,
			err,
		)
		return resp, err
	}

	d.logger.Printf(
		"[%s] Generate completed in %s",
		d.Name(),
		duration,
	)

	return resp, nil

}
