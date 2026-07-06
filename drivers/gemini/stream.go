package gemini

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/YashShekhawat/fusion/models"
)

type geminiStream struct {
	body   io.ReadCloser
	reader *bufio.Reader
}

func newGeminiStream(body io.ReadCloser) *geminiStream {
	return &geminiStream{
		body:   body,
		reader: bufio.NewReader(body),
	}
}

func (s *geminiStream) Close() error {
	return s.body.Close()
}

func (s *geminiStream) Recv() (models.StreamChunk, error) {

	for {

		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return models.StreamChunk{}, io.EOF
			}
			return models.StreamChunk{}, err
		}

		line = strings.TrimSpace(line)

		// Skip blank lines.
		if line == "" {
			continue
		}

		// Ignore event lines.
		if strings.HasPrefix(line, "event:") {
			continue
		}

		// We only process data lines.
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(
			strings.TrimPrefix(line, "data:"),
		)

		var resp streamResponse
		if err := json.Unmarshal([]byte(payload), &resp); err != nil {
			return models.StreamChunk{}, err
		}

		if len(resp.Candidates) == 0 {
			continue
		}

		candidate := resp.Candidates[0]

		if len(candidate.Content.Parts) == 0 {
			continue
		}

		var builder strings.Builder

		for _, part := range candidate.Content.Parts {
			if part.Text == "" {
				continue
			}
			builder.WriteString(part.Text)
		}

		text := strings.TrimSpace(builder.String())
		if text == "" {
			continue
		}

		fmt.Printf("Chunk received (%d chars)\n", len(text))
		return models.StreamChunk{
			Content: text,
		}, nil
	}
}
