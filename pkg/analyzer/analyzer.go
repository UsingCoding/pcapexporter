package analyzer

import (
	"context"
	stderrors "errors"
	stdlog "log"
	"os"
	"path/filepath"
	"pcapexporter/pkg/analyzer/crawler"
	"pcapexporter/pkg/analyzer/storage"
	"strings"
)

type Analyzer struct {
	Dir     string
	Workers uint
}

func (a Analyzer) Proceed(ctx context.Context) error {
	wp := NewWP(a.Workers)
	defer wp.Close()

	s, err := storage.NewCsvStorage("data.csv")
	if err != nil {
		return err
	}
	defer s.Close()

	err = filepath.Walk(a.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !strings.HasPrefix(filepath.Ext(path), ".pcap") {
			return nil
		}

		wp.Post(func() {
			handleErr := handleFile(ctx, path, s)

			stdlog.Printf("Complete with %s, err %s\n", path, handleErr)
		})

		return nil
	})

	wp.Wait()

	return err
}

func handleFile(ctx context.Context, p string, s storage.Storage) (err error) {
	res := make(chan crawler.Record)

	c := crawler.Crawler{
		Path:   p,
		Result: res,
	}

	go func() {
		processErr := c.Process(ctx)
		close(res)

		err = stderrors.Join(err, processErr)
	}()

	for rec := range res {
		storeRecord := s.StoreRecord(ctx, rec)
		err = stderrors.Join(err, storeRecord)
	}

	return err
}
