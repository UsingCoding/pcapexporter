package main

import (
	"github.com/urfave/cli/v2"
	pkganalyzer "pcapexporter/pkg/analyzer"
)

func analyzer() *cli.Command {
	return &cli.Command{
		Name:   "analyzer",
		Action: executeAnalyzer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to a pcap dir",
			},
			&cli.UintFlag{
				Name:  "workers",
				Usage: "Number of workers",
				Value: 20,
			},
		},
	}
}

func executeAnalyzer(c *cli.Context) error {
	return pkganalyzer.Analyzer{
		Dir:     c.String("path"),
		Workers: c.Uint("workers"),
	}.
		Proceed(c.Context)
}
