package versionfile

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var versionRE = regexp.MustCompile(`^go1($|\.)`)

type File struct {
	Version *Version
	Time    *Time
	Syntax  *FileSyntax
}

func Parse(file string, data []byte) (*File, error) {
	var version *Version
	var timeValue *Time
	first := true
	for line := range strings.Lines(string(data)) {
		line = strings.TrimRight(line, "\r\n")
		if first {
			first = false
			if !versionRE.MatchString(line) {
				return nil, fmt.Errorf("%s: first line does not match %s: %q", file, versionRE, line)
			}
			version = &Version{Version: line}
		} else {
			f := strings.Fields(line)
			if len(f) == 0 {
				continue
			}
			switch f[0] {
			case "time":
				if len(f) != 2 {
					return nil, fmt.Errorf("%s: expected two fields for 'time': %q", file, line)
				}
				t, err := time.ParseInLocation(time.RFC3339, f[1], time.UTC)
				if err != nil {
					return nil, fmt.Errorf("%s: bad time: %w", file, err)
				}
				timeValue = &Time{Time: t}
			default:
				return nil, fmt.Errorf("%s: unexpected line: %q", file, line)
			}
		}
	}
	return &File{Version: version, Time: timeValue}, nil
}

func (f *File) AddVersionStmt(version string) error {
	if !versionRE.MatchString(version) {
		return fmt.Errorf("invalid language version %q", version)
	}
	if f.Version != nil {
		f.Version.Version = version
	} else {
		f.Version = &Version{Version: version}
	}
	return nil
}

func (f *File) AddTimeStmt(t time.Time) error {
	if f.Time != nil {
		f.Time.Time = t
	} else {
		f.Time = &Time{Time: t}
	}
	return nil
}

func (f *File) DropTimeStmt() {
	f.Time = nil
}

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

type Version struct {
	Version string
}

type Time struct {
	Time time.Time
}

type FileSyntax struct {
	Name string // file path
}
