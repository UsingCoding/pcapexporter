package tshark

import (
	"context"
	"github.com/pkg/errors"
	stdlog "log"
	"os"
	"os/exec"
	"pcapexporter/pkg/common/fs"
)

type Converter struct{}

func (c Converter) ConvertToCsv(
	ctx context.Context,
	src, dst string,
) error {
	if fs.PathExists(dst) {
		stdlog.Println("SKIPPED, already exists:", dst)
		// don`t do job twice
		return nil
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.CommandContext(
		ctx,
		"tshark",
		"-r",
		src,
		"-E",
		"separator=,",
	)
	// write result to file directly
	cmd.Stdout = f
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		err = errors.Wrap(err, string(exitErr.Stderr))
	}
	return err

}
