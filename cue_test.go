package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/token"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
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

func TestRedact(t *testing.T) {
	v := ctx.CompileString(`
kind:       "Secret"
apiVersion: "v1"
metadata: {
	name:              "api-kt-config"
	namespace:         "rutang-beta"
	creationTimestamp: null
	ownerReferences: [{
		apiVersion: "bitnami.com/v1alpha1"
		kind:       "SealedSecret"
		name:       "api-kt-config"
		uid:        ""
		controller: true
	}]
}
data: hello: "aGVsbG8gd29ybGQ="
data: bar: "aGVsbG8gd29ybGQ="
`)
	redacted, err := DecodeSecret(&v)
	if err != nil {
		panic(err)
	}

	bytes, err := format.Node(redacted, format.Simplify())
	if err != nil {
		panic(err)
	}
	log.Println(string(bytes))
}

func TestAst(t *testing.T) {
	f := &ast.File{
		Decls: []ast.Decl{
			&ast.Package{Name: ast.NewIdent("foo")},
			&ast.EmbedDecl{
				Expr: &ast.BasicLit{
					Kind:     token.INT,
					ValuePos: token.NoSpace.Pos(),
					Value:    "1",
				},
			},
		},
	}
	node, err := format.Node(f)
	if err != nil {
		panic(err)
	}
	log.Println(string(node))
}

func TestInspectAst(t *testing.T) {
	v := ctx.CompileString(`
foo: bar: "baz"
`)
	s := v.Source().(*ast.File)
	dcls := s.Decls
	for _, dcl := range dcls {
		field := dcl.(*ast.Field)
		value := field.Value.(*ast.StructLit)
		log.Printf("field is %+v", field)
		log.Printf("value is %+v", value)
		for _, elt := range value.Elts {
			log.Printf("elt type %T", elt)
			log.Printf("elt is %+v", elt)

		}
	}

}

func TestBuildFile(t *testing.T) {
	v := ctx.CompileString(`
d1: 1
d2: 2
`)
	s := v.Syntax().(*ast.StructLit)
	f := BuildFile(&stringArray{"foo", "bar", "baz"}, "myPackage", s)
	bytes, err := format.Node(f, format.Simplify())
	if err != nil {
		panic(err)
	}
	log.Print(string(bytes))

	f = BuildFile(&stringArray{}, "myPackage", s)
	bytes, err = format.Node(f, format.Simplify())
	if err != nil {
		panic(err)
	}
	log.Print(string(bytes))
}

func TestBuildAstFile(t *testing.T) {
	f := &ast.File{
		Decls: []ast.Decl{
			&ast.Package{Name: ast.NewIdent("foo")},
			//&ast.NewIdent("foo: bar"),
			ast.NewStruct(
				&ast.Field{
					Label: ast.NewString("foo"),
					Value: ast.NewString("hello"),
				},
				&ast.Field{
					Label: ast.NewIdent("_bar"),
					Value: ast.NewStruct(
						&ast.Field{
							Label: ast.NewString("baz"),
							Value: ast.NewString("hello"),
						},
					),
				},
				&ast.Field{
					Label: ast.NewString("s1"),
					Value: ast.NewStruct(
						"foo", ast.NewString("f"),
						"baz", ast.NewString("b"),
					),
				},
			),
		},
	}
	v := ctx.BuildFile(f)
	fmt.Printf("%#v", v)

	bytes, _ := format.Node(f, format.Simplify())
	log.Printf("result is %v", string(bytes))
}

func TestInstance(t *testing.T) {
	var instance *cue.Value
	if i, err := getInstance(sampleCuePath(), "myapp"); err != nil {
		panic(err)
	} else {
		instance = i
	}

	log.Println(instance)
}

func TestGetTags(t *testing.T) {
	// todo
	// see cue/load/tags.go
}

func p(v interface{}) {
	log.Printf("value is %v\n", v)
}

func sampleCuePath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b) + "/sample"
}

// getInstance https://pkg.go.dev/cuelang.org/go/cue
func getInstance(dir string, root string) (*cue.Value, error) {
	config := load.Config{
		Context:    nil,
		ModuleRoot: dir,
		Module:     "",
		Package:    "",
		Dir:        dir,
		//Tags:        tags,
		TagVars:     nil,
		AllCUEFiles: false,
		Tests:       false,
		Tools:       false,
		DataFiles:   false,
		StdRoot:     "",
		ParseFile:   nil,
		Overlay:     nil,
		Stdin:       nil,
	}

	// eg. "./logging"
	args := []string{"./" + root}
	instances := load.Instances(args, &config)
	if l := len(instances); l != 1 {
		return nil, errors.New(fmt.Sprintf("can only evaluate exactly 1 cue instance, received %v", l))
	}
	instance := instances[0]

	log.Printf("tags are %v", instance.AllTags)

	var value cue.Value
	if v := ctx.BuildInstance(instance); v.Err() == nil {
		value = v
	} else {
		return nil, v.Err()
	}

	return &value, nil
}
