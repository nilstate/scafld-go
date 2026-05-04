package providers

import (
	"context"
	"io"
)

type LocalProvider struct {
	Messages []string
}

func (p LocalProvider) Invoke(ctx context.Context, w io.Writer) error {
	for _, msg := range p.Messages {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := io.WriteString(w, msg+"\n"); err != nil {
			return err
		}
	}
	return nil
}
