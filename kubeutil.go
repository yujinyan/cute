package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"encoding/base64"
)

func DecodeSecret(v *cue.Value) (*ast.StructLit, error) {
	data := v.LookupPath(cue.ParsePath("data"))

	fields, err := data.Fields()
	if err != nil {
		return nil, err
	}

	ret := *v

	for fields.Next() {
		value := fields.Value()

		selectors := value.Path().Selectors()
		selector := selectors[len(selectors)-1]

		var str string
		if s, err := value.String(); err != nil {
			return nil, err
		} else {
			str = s
		}

		var decoded []byte
		if d, err := base64.StdEncoding.DecodeString(str); err != nil {
			return nil, err
		} else {
			decoded = d
		}

		path := cue.MakePath(cue.Hid("_data", "_"), selector)
		ret = ret.FillPath(path, decoded)
	}

	syntax := ret.Syntax().(*ast.StructLit)
	var dataIdx int
	for i, elt := range syntax.Elts {
		elt := elt.(*ast.Field)
		label := elt.Label.(*ast.Ident)
		if label.Name == "data" {
			dataIdx = i
			break
		}
	}
	syntax.Elts = deleteAt(syntax.Elts, dataIdx)
	return syntax, nil
}

func BuildFile(labels *stringArray, pkgName string, s *ast.StructLit) *ast.File {

	f := &ast.File{}

	if pkgName != "" {
		f.Decls = append(f.Decls, &ast.Package{Name: ast.NewIdent(pkgName)})
	}

	if labels == nil || len(*labels) == 0 {
		f.Decls = append(f.Decls, s)
		return f
	}

	lastIdx := len(*labels) - 1
	cur := &ast.Field{}
	root := cur

	for i, label := range *labels {
		cur.Label = ast.NewIdent(label)
		if i < lastIdx {
			f := &ast.Field{}
			cur.Value = ast.NewStruct(f)
			cur = f
		}
	}

	cur.Value = s
	f.Decls = append(f.Decls, root)

	return f
}

func deleteAt(list []ast.Decl, idx int) []ast.Decl {
	lastIdx := len(list) - 1
	list[idx], list[lastIdx] = list[lastIdx], list[idx]
	return list[:lastIdx]
}
