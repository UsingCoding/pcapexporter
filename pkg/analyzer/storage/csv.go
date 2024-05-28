package storage

import (
	"context"
	"encoding/csv"
	stderrors "errors"
	"io"
	"os"
	"pcapexporter/pkg/analyzer/crawler"
	"strconv"
	"strings"
	"sync"
)

func NewCsvStorage(p string) (Storage, error) {
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(f)
	w := csvWriter{
		Writer: writer,
		Closer: f,
	}

	err = writer.Write(headers())
	if err != nil {
		return nil, stderrors.Join(err, w.Close())
	}

	return &csvStorage{
		writer: w,
		m:      sync.Mutex{},
	}, nil
}

type csvStorage struct {
	writer csvWriter
	m      sync.Mutex // protects write method
	num    int
}

func (c *csvStorage) StoreRecord(ctx context.Context, r crawler.Record) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.num++

	return c.writer.Write([]string{
		strconv.Itoa(c.num),
		strings.TrimSuffix(r.Timestamp.String(), " +0300 MSK"),
		r.File,
		strconv.Itoa(r.Seq),
		strconv.Itoa(r.RelativeID),
		r.Src,
		r.Dst,
		r.Data,
	})
}

func (c *csvStorage) Close() error {
	c.writer.Flush()
	return c.writer.Close()
}

type csvWriter struct {
	*csv.Writer
	io.Closer
}

func headers() []string {
	return []string{
		"record_id",
		"time",
		"file",
		"seq",
		"rel-id",
		"src",
		"dst",
		"data",
	}
}
