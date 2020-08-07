package main

import (
	"example/problem-2/urlshort"
	"fmt"
	"net/http"
)

func main() {
	mux := defaultMux()
	var dataformat string
	fmt.Println("Should we use yaml or json for url shortening?")
	fmt.Println("press 'y' for yaml and 'j' for json")
	fmt.Scanln(&dataformat)
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	//Build the YAMLHandler using the mapHandler as the
	yaml := `
  - path: /urlshort
    url: https://github.com/gophercises/urlshort
  - path: /urlshort-final
    url: https://github.com/gophercises/urlshort/tree/solution
`

	json := `[{"path": "/urlshort", "url": "https://github.com/gophercises/urlshort"},
	{"path": "/urlshort-final", "url": "https://github.com/gophercises/urlshort/tree/solution"}
	]`

	yamlhandler, err := urlshort.YAMLJSONHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	jsonhandler, err := urlshort.YAMLJSONHandler([]byte(json), mapHandler)
	_ = yamlhandler
	_ = jsonhandler

	if dataformat == "j" || dataformat == "y" {
		fmt.Println("Starting the server on :8080")
		if dataformat == "j" {
			http.ListenAndServe(":8080", jsonhandler)
		} else {
			http.ListenAndServe(":8080", jsonhandler)
		}

	} else {
		fmt.Println("you pressed an invalid key")
	}

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, razors")
}
