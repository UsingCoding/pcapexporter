package mysqlimport

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

func NewMySQL() (MySQL, error) {
	db, err := sqlx.Open("mysql", "grafana:1234@tcp(127.0.0.1:3306)/pcap")
	if err != nil {
		return MySQL{}, err
	}

	return MySQL{
		db: db,
	}, nil
}

type MySQL struct {
	db *sqlx.DB
}

type RecordSchema struct {
	Columns []string
}

type Records []string

func (m MySQL) Import(ctx context.Context, schema RecordSchema, recordsCh <-chan []Records) error {
	const insertSQLTpl = `INSERT INTO record (%s) VALUES %s`

	columns := makeColumns(schema.Columns)

	for records := range recordsCh {
		fmt.Println("Record", records)

		placeholders := makePlaceholders(
			uint(len(schema.Columns)),
			uint(len(records)),
		)

		insertSQL := fmt.Sprintf(insertSQLTpl, columns, placeholders)
		fmt.Println("SQL", insertSQL)

		_, err := m.db.ExecContext(ctx, insertSQL, emplace(records)...)
		if err != nil {
			return errors.Wrap(err, "failed to insert records")
		}
	}

	return nil
}

func (m MySQL) Close() error {
	return m.db.Close()
}

func makeColumns(columns []string) string {
	sanitized := make([]string, 0, len(columns))
	for _, column := range columns {
		sanitized = append(sanitized, fmt.Sprintf("`%s`", column))
	}

	return strings.Join(sanitized, ", ")
}

func makePlaceholders(columns, records uint) string {
	columnsPlaceholders := fmt.Sprintf(
		"(%s)",
		strings.Join(
			repeatV[string]("?", columns),
			",",
		),
	)

	return strings.Join(
		repeatV[string](columnsPlaceholders, records),
		", ",
	)
}

func repeatV[T any](v T, c uint) []T {
	res := make([]T, 0, c)
	for range c {
		res = append(res, v)
	}
	return res
}

func emplace(records []Records) []any {
	var res []any
	for _, record := range records {
		res = append(res, toAnySlice(record)...)
	}
	return res
}

func toAnySlice[T any](s []T) (res []any) {
	res = make([]any, 0, len(s))
	for _, v := range s {
		res = append(res, v)
	}
	return res
}
