package errors

import "errors"

var (
	// Registry Errors
	ErrDriverNotFound = errors.New("driver not found")
	ErrDuplicateDriver = errors.New("driver already registered")

	// Authentication Errors
	ErrUnauthorized = errors.New("unauthorized")

	// Rate Limiting
	ErrRateLimit = errors.New("rate limit exceeded")

	// Network / Provider Errors
	ErrTimeout = errors.New("request timed out")
	ErrProviderUnavailable = errors.New("provider unavailable")

	// Generic Errors
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidResponse = errors.New("invalid response")

	//streaming error
	ErrStreamingNotSupported = errors.New("streaming not supported")
)