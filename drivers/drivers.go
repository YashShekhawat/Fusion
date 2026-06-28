package drivers

import (
	"context"
	"github.com/YashShekhawat/fusion/models"
)

type Drivers interface {
	Name() string
	Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error)
}

