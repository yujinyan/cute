package cluster

import "encoding/base64"

secret: [Name=_]: {
	_data: [string]: string

	apiVersion: "v1"
	kind:       "Secret"
	metadata: {
		name:      Name
		namespace: #namespace
	}
	// https://pkg.go.dev/cuelang.org/go/pkg/encoding/base64
	if len(_data) > 0 {
		data: {for k, v in _data {"\(k)": base64.Encode(null, v)}}
	}
}

_seal: [Name=_]: {
	_data: [string]: string | bytes

	apiVersion: "v1"
	kind:       "Secret"
	metadata: {
		name:      Name
		namespace: #namespace
	}
	data: [string]: string

	// https://pkg.go.dev/cuelang.org/go/pkg/encoding/base64
	data: {for k, v in _data {"\(k)": base64.Encode(null, v)}}
}
