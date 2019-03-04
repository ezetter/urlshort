package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

var pathMap = make(map[string]string)

type PathMapping struct {
	Path string
	URL  string
}

func AddURL(path string, url string) {
	pathMap[path] = url
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matched := pathsToUrls[r.URL.Path]
		if matched != "" {
			fmt.Printf("From Map: Path = %v, matched=%v\n", r.URL.Path, matched)
			http.Redirect(w, r, matched, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r) // call original
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappingStructs []PathMapping
	err := yaml.Unmarshal(yml, &mappingStructs)
	if err != nil {
		return nil, err
	}

	for _, m := range mappingStructs {
		pathMap[m.Path] = m.URL
	}

	return MapHandler(pathMap, fallback), nil
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	matched := pathsToUrls[r.URL.Path]
	// 	if matched != "" {
	// 		fmt.Printf("From YAML: Path = %v, matched=%v\n", r.URL.Path, matched)
	// 		http.Redirect(w, r, matched, http.StatusSeeOther)
	// 	} else {
	// 		fallback.ServeHTTP(w, r) // call original
	// 	}
	// }), nil

}
