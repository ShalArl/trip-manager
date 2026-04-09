package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, fileName string, file io.Reader) error
	GetUrl(ctx context.Context, fileName string) (string, error)
}
