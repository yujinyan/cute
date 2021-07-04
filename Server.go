package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var ctx = cuecontext.New()
var dir string
var roots map[string]bool
var rootsList []string
var pathSplitter = func(c rune) bool {
	return c == '/'
}

func main() {
	//dir = "/home/yujinyan/code/rutang-cluster"
	roots = make(map[string]bool)
	dir = os.Getenv("CUE_DIR")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			roots[file.Name()] = true
		}
	}

	for dirname, _ := range roots {
		rootsList = append(rootsList, dirname)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.FieldsFunc(r.URL.Path, pathSplitter)
	//match := pathRe.FindStringSubmatch(r.URL.Path)
	root := pathComponents[0]
	log.Printf("root is %s", root)
	if !roots[root] {
		http.Error(w,
			fmt.Sprintf(`cannot find root "%s", available ones are %v`, root, rootsList),
			http.StatusBadRequest)
		return
	}
	var tags []string
	for key, values := range r.URL.Query() {
		if len(values) != 1 {
			http.Error(w,
				fmt.Sprintf(`tag must have exactly one value, key: "%s", values: %v"`, key, values),
				http.StatusBadRequest)
			return
		}
		tags = append(tags, fmt.Sprintf("%s=%s", key, values[0]))
	}
	// https://pkg.go.dev/cuelang.org/go/cue
	config := load.Config{
		Context:     nil,
		ModuleRoot:  dir,
		Module:      "",
		Package:     "",
		Dir:         dir,
		Tags:        tags,
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
		http.Error(w,
			fmt.Sprintf("can only evaluate exactly 1 cue instance, received %v", l),
			http.StatusBadRequest,
		)
	}

	instance := instances[0]

	values, err := ctx.BuildInstances(instances)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(values) < 1 {
		http.Error(w,
			fmt.Sprintf("root %s contains %v values, should contain only 1", root, len(values)),
			http.StatusBadRequest)
		return
	}

	value := values[0]

	var selectors []cue.Selector

	//log.Printf("package is %s\n", instances[0].ID())

	for _, seg := range pathComponents[1:] {
		if strings.HasPrefix(seg, "_") {
			log.Printf("add hidden\n")
			// https://github.com/cuelang/cue/issues/880
			// id: module/dir:package
			selectors = append(selectors, cue.Hid(seg, instance.ID()))
		} else if strings.HasPrefix(seg, "#") {
			// character `#` must url encode to `%23`
			selectors = append(selectors, cue.Def(seg))
		} else {
			selectors = append(selectors, cue.Str(seg))
		}
	}

	path := cue.MakePath(selectors...)
	log.Printf("path is %v\n", path)

	value = value.LookupPath(path)

	log.Printf("value is %v\n", value)

	result, err := yaml.Encode(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(result)
}
