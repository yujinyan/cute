package main

import (
	"path/filepath"
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

func TestGetDir(t *testing.T) {
	dir := sampleCuePath()

	testCases := []struct {
		paths  []string
		expect string
	}{
		{paths: []string{"myapp", "prod"}, expect: "myapp/prod"},
		{paths: []string{"myapp", "prod", "gibberish"}, expect: "myapp/prod"},
		{paths: []string{"myapp"}, expect: "myapp"},
	}

	for _, testCase := range testCases {
		idx := getDir(dir, testCase.paths)
		got := filepath.Join(testCase.paths[0:idx]...)
		if got != testCase.expect {
			t.Fatalf("expect %v, got %v", testCase.expect, got)
		}
	}
}
