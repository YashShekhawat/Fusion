package gemini

import (
	"fmt"
	"net/http"

	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
)

func mapStatusError(status int, body []byte) error {
	switch status {
	case http.StatusUnauthorized:
		return fmt.Errorf(
			"gemini: %w (status=%d): %s",
			fusionerrors.ErrUnauthorized,
			status,
			body,
		)

	case http.StatusRequestTimeout:
		return fmt.Errorf(
			"gemini: %w (status=%d): %s",
			fusionerrors.ErrTimeout,
			status,
			body,
		)

	case http.StatusTooManyRequests:
		return fmt.Errorf(
			"gemini: %w (status=%d): %s",
			fusionerrors.ErrRateLimit,
			status,
			body,
		)

	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return fmt.Errorf(
			"gemini: %w (status=%d): %s",
			fusionerrors.ErrProviderUnavailable,
			status,
			body,
		)

	default:
		return fmt.Errorf(
			"gemini: %w (status=%d): %s",
			fusionerrors.ErrInvalidResponse,
			status,
			body,
		)
	}
}
