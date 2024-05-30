package fs

import "os"

func PathExists(p string) bool {
	_, err := os.Stat(p)

	return !os.IsNotExist(err)
}
