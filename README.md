# Cute

A small [cue-lang](https://github.com/cue-lang/cue) utility.

## Features

### HTTP server for cue files

https://cuelang.org/docs/concepts/packages/#instances

All cue files with the same package can be evaluated within the context of a certain directory. Within this context,
only the files belonging to that package in that directory and its ancestor directories within the module are combined.
This is called an *instance* of a package.

The module root is marked by the presence of a `cue.mod` directory, initialized by `cue mod init <mymodule>`.

Running the `cute` command at a cue module root starts an HTTP server. The server evaluates an instance of cue
corresponding to the request path and returns the object in YAML, ready to be piped into `kubectl`. The cue module root
can also be specified in the `CUE_DIR` environment variable.

```shell
cd sample
cute
```

Example:

```shell
$ cd sample
$ tree .
.
├── config.cue
├── cue.mod
│   ├── module.cue
│   ├── pkg
│   └── usr
├── myapp
│   └── config.cue
└── namespace.cue

$ cue export ./myapp/ -e  'secret["my-app-config"]' --out yaml
apiVersion: v1
kind: Secret
data:
  foo: SGVsbG8gV29ybGQh
metadata:
  name: my-app-config
  namespace: default
  
$ cute &

$ curl http://localhost:8080/myapp
secret:
  my-app-config:
    apiVersion: v1
    kind: Secret
    data:
      foo: SGVsbG8gV29ybGQh
    metadata:
      name: my-app-config
      namespace: default

$ curl http://localhost:8080/myapp/secret/my-app-config
apiVersion: v1
kind: Secret
data:
  foo: SGVsbG8gV29ybGQh
metadata:
  name: my-app-config
  namespace: default
```

Use query parameters to inject tagged values (`-t` flags).

```sh
$ curl "http://localhost:8080/myapp/secret/my-app-config?namespace=prod"
apiVersion: v1
kind: Secret
data:
  foo: SGVsbG8gV29ybGQh
metadata:
  name: my-app-config
  namespace: prod
```

### Base64 decode Kubernetes Secret data

```shell
$ cd sample/myapp
$ cue eval -e 'secret["my-app-config"]' | cute
{
        apiVersion: "v1"
        kind:       "Secret"
        metadata: {
                name:      "my-app-config"
                namespace: "default"
        }
        _data: {
                foo: 'Hello World!'
        }
}
```

Path and cue package can be specified with `-l` and `-p` flags.

```shell
$ cue eval -e 'secret["my-app-config"]' | cute -l secret -l "my-app-config" -p "mycluster"
package mycluster

secret: "my-app-config": {
        apiVersion: "v1"
        kind:       "Secret"
        metadata: {
                name:      "my-app-config"
                namespace: "default"
        }
        _data: {
                foo: 'Hello World!'
        }
}
```


