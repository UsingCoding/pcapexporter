package provider

import (
	"context"
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	stdlog "log"
	"pcapexporter/pkg/common"
	"pcapexporter/pkg/mysqlimport"
)

type CSVProvider struct {
	Reader   io.ReadCloser
	PackSize int
}

func (p CSVProvider) Fetch(ctx context.Context) (
	mysqlimport.RecordSchema,
	<-chan []mysqlimport.Records,
	error,
) {
	r := csv.NewReader(p.Reader)

	header, err := r.Read()
	if err != nil {
		return mysqlimport.RecordSchema{}, nil, errors.Wrapf(err, "failed to read CSV header")
	}

	schema := mysqlimport.RecordSchema{
		Columns: header,
	}

	records := common.NewPackChan(
		50,
		make(chan []mysqlimport.Records, 1),
	)

	go func() {
		for {
			record, err := r.Read()
			if err != nil {
				records.Close()
				if errors.Is(err, io.EOF) {
					break
				}
				stdlog.Println("CSV err", err)
				break
			}

			records.Send(record)
		}
	}()

	return schema, records.Recv(), nil
}
