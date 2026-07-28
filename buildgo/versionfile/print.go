package versionfile

import (
	"bytes"
	"fmt"
	"time"
)

func (f *File) Format() ([]byte, error) {
	b := &bytes.Buffer{}
	if f.Version == nil {
		return nil, fmt.Errorf("f.Version is nil")
	}
	fmt.Fprintf(b, "%s\n", f.Version.Version)
	if f.Time != nil {
		fmt.Fprintf(b, "%s %s\n", "time", f.Time.Time.UTC().Format(time.RFC3339))
	}
	return b.Bytes(), nil
}
