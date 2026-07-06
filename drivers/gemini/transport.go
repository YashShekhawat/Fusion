package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
)

// sendRequest is scoped to GeminiDriver and handles
// sending JSON payloads to Gemini and decoding the response.
func (d *GeminiDriver) sendRequest(ctx context.Context, endpoint string, payload any, out any) error {
	// Marshal the payload into JSON.
	jsonReq, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(
			"gemini: marshal request: %w",
			fusionerrors.ErrInvalidRequest,
		)
	}

	// Construct the HTTP request.
	url := d.baseURL + endpoint
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(jsonReq),
	)
	if err != nil {
		return fmt.Errorf(
			"gemini: create request: %w",
			fusionerrors.ErrInvalidRequest,
		)
	}

	httpReq.Header.Set("x-goog-api-key", d.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute the request.
	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf(
				"gemini: send request: %w",
				fusionerrors.ErrTimeout,
			)
		}

		return fmt.Errorf(
			"gemini: send request: %w",
			fusionerrors.ErrProviderUnavailable,
		)
	}
	defer resp.Body.Close()

	// Validate the HTTP status.
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"gemini: read error response: %w",
				fusionerrors.ErrInvalidResponse,
			)
		}

		return mapStatusError(resp.StatusCode, body)
	}

	// Decode the response.
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf(
			"gemini: decode response: %w",
			fusionerrors.ErrInvalidResponse,
		)
	}

	return nil
}

func (d *GeminiDriver) sendStreamRequest(ctx context.Context, endpoint string, payload any) (io.ReadCloser, error) {

	// Marshal the payload.
	jsonReq, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf(
			"gemini: marshal request: %w",
			fusionerrors.ErrInvalidRequest,
		)
	}

	// Create request.
	url := d.baseURL + endpoint

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(jsonReq),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gemini: create request: %w",
			fusionerrors.ErrInvalidRequest,
		)
	}

	httpReq.Header.Set("x-goog-api-key", d.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request.
	resp, err := d.httpClient.Do(httpReq)
	if err != nil {

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf(
				"gemini: send request: %w",
				fusionerrors.ErrTimeout,
			)
		}

		return nil, fmt.Errorf(
			"gemini: send request: %w",
			fusionerrors.ErrProviderUnavailable,
		)
	}

	// Validate status.
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			return nil, fmt.Errorf(
				"gemini: read error response: %w",
				fusionerrors.ErrInvalidResponse,
			)
		}

		return nil, mapStatusError(resp.StatusCode, body)
	}

	// Success.
	// IMPORTANT:
	// we are not closing the body as stream implementation owns it.
	return resp.Body, nil
}
