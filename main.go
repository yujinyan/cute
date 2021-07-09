package main

import (
	"cuelang.org/go/cue/cuecontext"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var ctx = cuecontext.New()

func main() {
	if hasStdIn() {
		var packageArg = flag.String("p", "", "package name for generated cue file")
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
		if *packageArg != "" {
			fmt.Printf("package %s\n\n", *packageArg)
		}
		fmt.Printf("%#v", redacted)
		return
	}

	ServeCueFiles()
}

func hasStdIn() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
