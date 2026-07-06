package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendRequestMapsHTTPStatusToFusionErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       error
	}{
		{name: "unauthorized", statusCode: http.StatusUnauthorized, want: fusionerrors.ErrUnauthorized},
		{name: "rate limit", statusCode: http.StatusTooManyRequests, want: fusionerrors.ErrRateLimit},
		{name: "timeout", statusCode: http.StatusRequestTimeout, want: fusionerrors.ErrTimeout},
		{name: "provider unavailable", statusCode: http.StatusBadGateway, want: fusionerrors.ErrProviderUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &OpenAIDriver{
				apiKey:  "test-key",
				baseURL: "https://example.com",
				httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.statusCode,
						Status:     fmt.Sprintf("%d %s", tt.statusCode, http.StatusText(tt.statusCode)),
						Body:       io.NopCloser(strings.NewReader(`{"error":"boom"}`)),
						Header:     make(http.Header),
						Request:    req,
					}, nil
				})},
			}

			err := d.sendRequest(context.Background(), "/responses", map[string]string{"foo": "bar"}, &map[string]any{})
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected error to wrap %v, got %v", tt.want, err)
			}
		})
	}
}
