package main

import (
	"strings"
	"testing"
)

var splitFn = func(c rune) bool {
	return c == '/'
}

func TestParsePath(t *testing.T) {
	path := "/foo/bar"
	components := strings.FieldsFunc(path, splitFn)

	if len(components) != 2 {
		t.Fatalf("wrong size")
	}

	if components[0] != "foo" || components[1] != "bar" {
		t.Fatalf("wrong content")
	}
}
