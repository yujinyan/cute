package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var dir string
var roots map[string]bool
var rootsList []string
var pathSplitter = func(c rune) bool {
	return c == '/'
}

func ServeCueFiles() {
	roots = make(map[string]bool)
	var dir string
	if s := os.Getenv("CUE_DIR"); s != "" {
		dir = s
	} else {
		s, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}
		dir = s
	}

	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	var files []fs.FileInfo
	if f, err := ioutil.ReadDir(dir); err == nil {
		files = f
	} else {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			roots[file.Name()] = true
		}
	}

	for dirname := range roots {
		rootsList = append(rootsList, dirname)
	}

	http.HandleFunc("/", handler)
	log.Printf("cue server started for dir %s\n", dir)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.FieldsFunc(r.URL.Path, pathSplitter)
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
	cut := getDir(dir, pathComponents)
	// https://pkg.go.dev/cuelang.org/go/cue
	config := load.Config{
		Context:     nil,
		ModuleRoot:  dir,
		Module:      "",
		Package:     "",
		Dir:         filepath.Join(pathComponents[0:cut]...),
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

	instances := load.Instances(nil, &config)
	if l := len(instances); l != 1 {
		http.Error(w,
			fmt.Sprintf("can only evaluate exactly 1 cue instance, received %v", l),
			http.StatusBadRequest,
		)
	}
	instance := instances[0]

	var value cue.Value
	if v := ctx.BuildInstance(instance); v.Err() == nil {
		value = v
	} else {
		panic(v.Err())
	}

	var selectors []cue.Selector
	for _, seg := range pathComponents[cut:] {
		if strings.HasPrefix(seg, "_") {
			// https://github.com/cuelang/cue/issues/880
			// id format: module/dir:package
			selectors = append(selectors, cue.Hid(seg, instance.ID()))
		} else if strings.HasPrefix(seg, "#") {
			// character `#` must url encode to `%23`
			selectors = append(selectors, cue.Def(seg))
		} else {
			selectors = append(selectors, cue.Str(seg))
		}
	}
	path := cue.MakePath(selectors...)
	value = value.LookupPath(path)

	var result []byte
	var resultErr error
	if list, err := value.List(); err != nil {
		result, resultErr = yaml.Encode(value)
	} else {
		result, resultErr = yaml.EncodeStream(list)
	}
	if resultErr != nil {
		http.Error(w, resultErr.Error(), http.StatusBadRequest)
		return
	}
	_, _ = w.Write(result)
}

// navigate down file path according to pathComponents
func getDir(root string, pathComponents []string) int {
	var i int = 1
	for i <= len(pathComponents) {
		comps := append([]string{root}, pathComponents[0:i]...)
		cur := filepath.Join(comps...)
		if _, err := os.Stat(cur); os.IsNotExist(err) {
			break
		} else if err == nil {
			i++
		} else {
			panic(err)
		}
	}
	return i - 1
}
