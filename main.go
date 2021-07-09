package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
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
		fmt.Printf("%#v", redacted)
		return
	}

	ServeCueFiles()
}

func hasStdIn() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
