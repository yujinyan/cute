package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"log"
	"testing"
)

func TestDecodeSecret(t *testing.T) {
	v := ctx.CompileString(`
kind:       "Secret"
apiVersion: "v1"
metadata: {
	name:              "my-secret"
	namespace:         "default"
	creationTimestamp: null
	ownerReferences: [{
		apiVersion: "bitnami.com/v1alpha1"
		kind:       "SealedSecret"
		name:       "my-secret"
		uid:        ""
		controller: true
	}]
}
data: foo: "aGVsbG8gd29ybGQ="
data: bar: "aGVsbG8gd29ybGQ="
`)
	decoded, err := DecodeSecret(&v)
	if err != nil {
		panic(err)
	}

	file := BuildFile(nil, "", decoded)

	v = ctx.BuildFile(file)
	path := cue.MakePath(cue.Hid("_data", "_"), cue.Str("foo"))
	got := v.LookupPath(path)
	expect := "hello world"
	if bytes, _ := got.Bytes(); "hello world" != string(bytes) {
		t.Fatalf("expect %v, got %v", expect, got)
	}

	bytes, err := format.Node(decoded, format.Simplify())
	if err != nil {
		panic(err)
	}
	log.Println(string(bytes))
}

func TestBuildFile(t *testing.T) {
	v := ctx.CompileString(`
d1: 1
d2: 2
`)
	s := v.Syntax().(*ast.StructLit)
	testCases := []struct {
		labels  *stringArray
		pkgName string
		expect  string
	}{
		{&stringArray{"foo"}, "myPackage", `package myPackage

foo: {
	d1: 1
	d2: 2
}
`},
		{&stringArray{"foo", "bar", "baz"}, "myPackage", `package myPackage

foo: bar: baz: {
	d1: 1
	d2: 2
}
`},
		{&stringArray{}, "myPackage", `package myPackage

{
	d1: 1
	d2: 2
}
`},
		{nil, "myPackage", `package myPackage

{
	d1: 1
	d2: 2
}
`},
	}

	for _, testCase := range testCases {
		f := BuildFile(testCase.labels, testCase.pkgName, s)
		bytes, _ := format.Node(f, format.Simplify())
		if got := string(bytes); got != testCase.expect {
			t.Fatalf("expect %v, got %v", testCase.expect, got)
		}
	}
}
