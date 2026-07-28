package osutil

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"time"
)

func PollStat(ctx context.Context, name string, d time.Duration) (fs.FileInfo, error) {
	for {
		info, err := os.Stat(name)
		if err == nil {
			return info, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d):
		}
	}
}
