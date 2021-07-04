package main

import (
	"log"
	"strings"
	"testing"
)

var splitFn = func(c rune) bool {
	return c == '/'
}

func TestParsePath(t *testing.T) {
	path := "/api-backend/objects"
	components := strings.FieldsFunc(path, splitFn)
	log.Println(components)
}
