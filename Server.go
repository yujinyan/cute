package main

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var pathRe = regexp.MustCompile(`^/(.+)`)
var ctx = cuecontext.New()

var dir string
var roots map[string]bool
var rootsList []string

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
	match := pathRe.FindStringSubmatch(r.URL.Path)
	path := match[1]
	log.Printf("path is %s", path)
	if !roots[path] {
		http.Error(w,
			fmt.Sprintf(`cannot find root path "%s", available ones are %v`, path, rootsList),
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
	config := load.Config{
		Context:     nil,
		ModuleRoot:  dir,
		Module:      "",
		Package:     "",
		Dir:         dir,
		Tags:        tags,
		TagVars:     nil,
		AllCUEFiles: false,
		BuildTags:   nil,
		Tests:       false,
		Tools:       false,
		DataFiles:   false,
		StdRoot:     "",
		ParseFile:   nil,
		Overlay:     nil,
		Stdin:       nil,
	}

	// eg. "./logging"
	args := []string{"./" + path}
	instances := load.Instances(args, &config)

	values, err := ctx.BuildInstances(instances)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, value := range values {
		result, err := yaml.Encode(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(result)
	}
}
