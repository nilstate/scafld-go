package providers

import (
	"context"
	"strings"

	"github.com/nilstate/scafld-go/internal/core/review"
)

type LocalProvider struct {
	Messages []string
}

func (p LocalProvider) Invoke(ctx context.Context, taskID string) (review.Packet, error) {
	var lines []string
	for _, msg := range p.Messages {
		if err := ctx.Err(); err != nil {
			return review.Packet{}, err
		}
		lines = append(lines, msg)
	}
	if len(lines) == 0 {
		lines = []string{`{"type":"verdict","verdict":"pass"}`}
	}
	return review.ParseNDJSON(strings.Join(lines, "\n") + "\n")
}
