package drivers

import (
	"context"

	"github.com/YashShekhawat/fusion/models"
)

// Stream represents a streaming text generation response.
//
// Recv blocks until the next chunk is available.
// When the stream finishes, Recv returns io.EOF.
// Close releases any underlying resources.
type Stream interface {
	Recv() (models.StreamChunk, error)
	Close() error
}

type StreamDriver interface {
	Driver

	GenerateStream(ctx context.Context, req models.GenerateRequest) (Stream, error)
}
