package analyzer

import (
	"context"
	stderrors "errors"
	"fmt"
	stdlog "log"
	"os"
	"path/filepath"
	"pcapexporter/pkg/analyzer/crawler"
	"pcapexporter/pkg/analyzer/storage"
	"pcapexporter/pkg/analyzer/tshark"
	"strings"
	"time"
)

type Analyzer struct {
	Dir     string
	Result  string
	Workers uint
}

func (a Analyzer) Proceed(ctx context.Context) error {
	wp := NewWP(a.Workers)
	defer wp.Close()

	s, err := storage.NewCsvStorage(a.Result)
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
			start := time.Now()

			handleErr := handleFile(ctx, path, s)

			msg := fmt.Sprintf("Complete with %s, took: %v", path, time.Since(start))
			if handleErr != nil {
				msg += fmt.Sprintf(", err: %s", handleErr)
			}

			stdlog.Println(msg)
		})

		return nil
	})

	wp.Wait()

	return err
}

func handleFile(ctx context.Context, p string, s storage.Storage) (err error) {
	csvPath := fmt.Sprintf("%s.csv", p)

	err = tshark.Converter{}.ConvertToCsv(ctx, p, csvPath)
	if err != nil {
		return err
	}

	res := make(chan crawler.Record)

	c := crawler.Crawler{
		Path:    p,
		CsvPath: csvPath,
		Result:  res,
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
