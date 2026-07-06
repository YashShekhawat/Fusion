package middleware

import (
	"context"
	"io"
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

type loggingStreamDriver struct {
	*loggingDriver
	streamNext drivers.StreamDriver
}

type loggingStream struct {
	next   drivers.Stream
	logger Logger
	name   string
	start  time.Time
}

func Logging(logger Logger) Middleware {
	return func(d drivers.Driver) drivers.Driver {

		base := &loggingDriver{
			next:   d,
			logger: logger,
		}

		// Preserve streaming capability.
		if streamDriver, ok := d.(drivers.StreamDriver); ok {
			return &loggingStreamDriver{
				loggingDriver: base,
				streamNext:    streamDriver,
			}
		}

		return base
	}
}

func (d *loggingDriver) Name() string {
	return d.next.Name()
}

func (d *loggingDriver) Generate(ctx context.Context, req models.GenerateRequest) (models.GenerateResponse, error) {

	start := time.Now()

	d.logger.Printf("[%s] Generate started", d.Name())

	resp, err := d.next.Generate(ctx, req)

	if err != nil {
		d.logger.Printf(
			"[%s] Generate failed after %s: %v",
			d.Name(),
			time.Since(start),
			err,
		)
		return resp, err
	}

	d.logger.Printf(
		"[%s] Generate completed in %s",
		d.Name(),
		time.Since(start),
	)

	return resp, nil
}

func (d *loggingStreamDriver) GenerateStream(ctx context.Context, req models.GenerateRequest) (drivers.Stream, error) {

	start := time.Now()

	d.logger.Printf("[%s] Stream started", d.Name())

	stream, err := d.streamNext.GenerateStream(ctx, req)
	if err != nil {

		d.logger.Printf(
			"[%s] Stream failed after %s: %v",
			d.Name(),
			time.Since(start),
			err,
		)

		return nil, err
	}

	return &loggingStream{
		next:   stream,
		logger: d.logger,
		name:   d.Name(),
		start:  start,
	}, nil
}

func (s *loggingStream) Recv() (models.StreamChunk, error) {

	chunk, err := s.next.Recv()

	switch err {

	case nil:
		return chunk, nil

	case io.EOF:
		s.logger.Printf(
			"[%s] Stream completed in %s",
			s.name,
			time.Since(s.start),
		)
		return chunk, err

	default:
		s.logger.Printf(
			"[%s] Stream failed after %s: %v",
			s.name,
			time.Since(s.start),
			err,
		)
		return chunk, err
	}
}

func (s *loggingStream) Close() error {
	return s.next.Close()
}
