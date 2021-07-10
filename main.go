package main

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"flag"
	"io/ioutil"
	"os"
)

type stringArray []string

func (i *stringArray) String() string {
	return "string array"
}

func (i *stringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var ctx = cuecontext.New()
var labelsArg stringArray

func main() {
	if hasStdIn() {
		var packageArg = flag.String("p", "", "package name for generated cue file")
		flag.Var(&labelsArg, "l", "label path")
		flag.Parse()

		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		v := ctx.CompileBytes(bytes)
		redacted, err := DecodeSecret(&v)
		if err != nil {
			panic(err)
		}

		f := BuildFile(&labelsArg, *packageArg, redacted)

		result, err := format.Node(f, format.Simplify())
		if err != nil {
			panic(err)
		}

		os.Stdout.Write(result)
		return
	}

	ServeCueFiles()
}

func hasStdIn() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
