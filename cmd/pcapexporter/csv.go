package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	stdlog "log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func csv() *cli.Command {
	return &cli.Command{
		Name:   "csv",
		Action: executeCSV,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to a pcap dir",
			},
		},
	}
}

func executeCSV(c *cli.Context) error {
	ctx := c.Context
	dir := c.String("path")

	stdlog.Println("Start")

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		if !strings.HasPrefix(filepath.Ext(path), ".pcap") {
			return nil
		}

		resPath := fmt.Sprintf("%s.csv", path)

		if fileExists(resPath) {
			stdlog.Println("SKIPPED, already exists:", path)
			return nil
		}

		f, err := os.Create(resPath)
		if err != nil {
			return err
		}
		defer f.Close()

		start := time.Now()

		cmd := exec.CommandContext(
			ctx,
			"tshark",
			"-r",
			path,
			"-E",
			"separator=,",
		)
		// write result to file directly
		cmd.Stdout = f
		cmd.Stderr = os.Stderr

		cmdErr := cmd.Run()
		stdlog.Printf("Completed with %s, took: %v\n", path, time.Since(start))

		if cmdErr != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				err = errors.Wrap(err, string(exitErr.Stderr))
			}
			stdlog.Println("CMD Err:", cmdErr)
		}
		return nil
	})
}

func fileExists(p string) bool {
	_, err := os.Stat(p)

	return !os.IsNotExist(err)
}
