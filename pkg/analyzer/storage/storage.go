package storage

import (
	"context"
	"io"
	"pcapexporter/pkg/analyzer/crawler"
)

type Storage interface {
	StoreRecord(ctx context.Context, r crawler.Record) error
	io.Closer
}
