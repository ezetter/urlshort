package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ezetter/urlshort"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func loadYaml(fileName string) []byte {
	dat, err := ioutil.ReadFile(fileName)
	check(err)

	return dat
}

var pathsToUrls = map[string]string{
	"/urlshort-godocy": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godocz":     "https://godoc.org/gopkg.in/yaml.v2",
}

func main() {
	mux := defaultMux()

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(loadYaml("mappings.yaml"), mapHandler)
	check(err)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", catchall)
	mux.HandleFunc("/addUrl", addURL)
	return mux
}

func addURL(w http.ResponseWriter, r *http.Request) {
	pairs := r.URL.Query()
	for k, v := range pairs {
		if strings.HasPrefix(k, "/") && k != "/addUrl" {
			pathsToUrls[k] = v[0]
			fmt.Fprintf(w, "Added url %s -> %s\n", k, v[0])
		} else {
			fmt.Fprintf(w, "Query string is malformed.")
		}
	}
}

func catchall(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	fmt.Printf("No matches.")
}
