package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ramin0/live/go/urlshort/urlshort"
)

func main() {
	flagYamlFilename := flag.String("yml", "urls.yaml", "Path to YAML file containing path/url mappings")
	flagJSONFilename := flag.String("json", "urls.json", "Path to JSON file containing path/url mappings")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{ // make(map[string]string)
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
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
	ymlFile, err := os.Open(*flagYamlFilename)
	if err != nil {
		fmt.Printf("failed to open %q: %v\n", *flagYamlFilename, err)
		return
	}
	defer ymlFile.Close()
	yaml, err := ioutil.ReadAll(ymlFile)
	if err != nil {
		fmt.Printf("failed to read %q: %v\n", *flagYamlFilename, err)
		return
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonFile, err := os.Open(*flagJSONFilename)
	if err != nil {
		fmt.Printf("failed to open %q: %v\n", *flagJSONFilename, err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("failed to read %q: %v\n", *flagJSONFilename, err)
		return
	}
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonData), yamlHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
