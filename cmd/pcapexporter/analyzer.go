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
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Path to a pcap dir",
			},
			&cli.UintFlag{
				Name:    "workers",
				Aliases: []string{"w"},
				Usage:   "Number of workers",
				Value:   20,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Path to result csv file",
			},
		},
	}
}

func executeAnalyzer(c *cli.Context) error {
	return pkganalyzer.Analyzer{
		Dir:     c.String("path"),
		Workers: c.Uint("workers"),
		Result:  c.String("output"),
	}.
		Proceed(c.Context)
}
