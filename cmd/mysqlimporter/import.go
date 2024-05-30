package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"pcapexporter/pkg/mysqlimport"
	"pcapexporter/pkg/mysqlimport/provider"
)

func importAction(ctx *cli.Context) error {
	db, err := mysqlimport.NewMySQL()
	if err != nil {
		return err
	}
	defer db.Close()

	csvProvider := provider.CSVProvider{Reader: os.Stdin}

	schema, records, err := csvProvider.Fetch(ctx.Context)
	if err != nil {
		return err
	}

	return db.Import(ctx.Context, schema, records)
}
