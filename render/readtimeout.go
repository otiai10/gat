// +build !go1.10

package render

import (
	"os"
	"time"
)

func readTimeout(f *os.File, t time.Time) error {
	return nil
}
