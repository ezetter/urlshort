package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godocy": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godocz":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	yamlHandler, err := urlshort.YAMLHandler(loadYaml("mappings.yaml"), mapHandler)
	check(err)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
