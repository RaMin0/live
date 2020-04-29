package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// fmt.Printf("Path: %q\n", path)
		longURL, ok := pathsToUrls[path]
		if !ok {
			// couldn't find the request's path in the map
			fallback.ServeHTTP(w, r)
			return
		}
		// otherwise, redirect to longURL
		// http.Redirect(w, r, longURL, http.StatusMovedPermanently)
		fmt.Fprintf(w, "Redirecting to %s\n", longURL)
	}
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

type yamlData struct {
	Path string
	URL  string
}

func YAMLHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the YAML raw data into the yamlData struct
	// var parsedYamlData []map[string]string
	var parsedYamlData []yamlData
	if err := yaml.Unmarshal([]byte(data), &parsedYamlData); err != nil {
		return nil, err
	}
	// fmt.Printf("%+v\n", parsedYamlData)
	// build pathsToUrls from the parsed YAML data
	pathsToUrls := map[string]string{}
	for _, yamlEntry := range parsedYamlData {
		pathsToUrls[yamlEntry.Path] = yamlEntry.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the YAML raw data into the yamlData struct
	// var parsedYamlData []map[string]string
	// this is very similar to the YAML handler, we just switch yaml with json!
	var parsedYamlData []yamlData
	if err := json.Unmarshal([]byte(data), &parsedYamlData); err != nil {
		return nil, err
	}
	// fmt.Printf("%+v\n", parsedYamlData)
	// build pathsToUrls from the parsed YAML data
	pathsToUrls := map[string]string{}
	for _, yamlEntry := range parsedYamlData {
		pathsToUrls[yamlEntry.Path] = yamlEntry.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}
