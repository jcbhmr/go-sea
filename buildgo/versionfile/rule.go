package versionfile

import "time"

type File struct {
	Version *Version
	Time    *Time
	Syntax  *FileSyntax
}

type Version struct {
	Version string
}

type Time struct {
	Time time.Time
}