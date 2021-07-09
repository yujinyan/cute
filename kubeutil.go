package main

import (
	"cuelang.org/go/cue"
	"encoding/base64"
)

func DecodeSecret(v *cue.Value) (*cue.Value, error) {
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

		path := cue.MakePath(cue.Def("data"), selector)
		ret = ret.FillPath(path, decoded)
	}
	return &ret, nil
}
