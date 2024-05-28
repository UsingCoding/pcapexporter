package storage

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"pcapexporter/pkg/analyzer/crawler"
)

func NewMySQLStorage() (Storage, error) {
	db, err := sqlx.Open("mysql", "grafana:1234@tcp(127.0.0.1:3306)/pcap")
	if err != nil {
		return nil, err
	}

	return storage{
		db: db,
	}, nil
}

type storage struct {
	db *sqlx.DB
}

func (s storage) StoreRecord(ctx context.Context, r crawler.Record) error {
	const insertSQL = `
		INSERT INTO record (
			time,
			file,
			seq,
			src,
			dst,
			data
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, insertSQL, r.Timestamp, r.File, r.Seq, r.Src, r.Dst, r.Data)
	return errors.Wrap(err, "failed to store record")
}

func (s storage) Close() error {
	return s.db.Close()
}
