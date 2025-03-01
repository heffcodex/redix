package redix

import "strings"

const (
	KeyDelimiter = ":"
)

type Namespace string

func (ns Namespace) String() string {
	return string(ns)
}

func (ns Namespace) Append(parts ...string) Namespace {
	var slice []string

	if ns != "" {
		slice = []string{string(ns)}
	}

	return Namespace(strings.Join(append(slice, parts...), KeyDelimiter))
}
