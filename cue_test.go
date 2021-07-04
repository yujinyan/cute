package main

import (
	"cuelang.org/go/cue"
	"log"
	"testing"
)

func TestLookupHidden(t *testing.T) {
	v := ctx.CompileString(`
package test
_v: 1
v: 1
`)

	selector := cue.Hid("_v", "_")
	path := cue.MakePath(selector)
	p(v.LookupPath(path))
}

func TestLookupNormalField(t *testing.T) {
	v := ctx.CompileString(`
package test
_v: 1
v: 1
`)

	selector := cue.Str("v")
	path := cue.MakePath(selector)
	p(v.LookupPath(path))
}

func TestLookupHidWithoutPackage(t *testing.T) {
	v := ctx.CompileString(`
_v: 1
v: 1
`)

	selector := cue.Hid("_v", "_")
	path := cue.MakePath(selector)
	p(v.LookupPath(path))
}

func p(v interface{}) {
	log.Printf("value is %v\n", v)
}
