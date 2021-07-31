package main

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"flag"
	"fmt"
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
var (
	version = "dev"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	var packageArg = flag.String("p", "", "package name for generated cue file")
	var printVersion = flag.Bool("version", false, "print version")
	flag.Var(&labelsArg, "l", "label path")
	flag.Parse()

	if *printVersion {
		fmt.Printf("cute version: %s, built by %s on %s.\n", version, builtBy, date)
		return
	}

	if hasStdIn() {
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
