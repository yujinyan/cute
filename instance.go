package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"fmt"
)

func GetInstance(dir string) (*cue.Value, error) {
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
	args := []string{"./" + "api-kt"}
	instances := load.Instances(args, &config)
	if l := len(instances); l != 1 {
		return nil, errors.New(fmt.Sprintf("can only evaluate exactly 1 cue instance, received %v", l))
	}
	instance := instances[0]

	//instance.AddSyntax()

	var value cue.Value
	if v := ctx.BuildInstance(instance); v.Err() == nil {
		value = v
	} else {
		panic(v.Err())
	}

	return &value, nil
}
