package urlshort

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	yaml "gopkg.in/yaml.v2"
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
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
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
func YAMLJSONHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYamlJson(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := converttomap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func converttomap(pathUrls []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYamlJson(data []byte) ([]pathURL, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func DbHandler(database string, user string, password string, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := connectdatabase(database, user, password)
	if err != nil {
		return nil, err
	}
	return MapHandler(paths, fallback), nil
}

func connectdatabase(database string, user string, password string) (map[string]string, error) {

	db, err := sql.Open("mysql", user+":"+password+"@tcp(127.0.0.1:3306)/"+database)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM paths")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	m := make(map[string]string)
	for rows.Next() {
		var path pathURL

		err = rows.Scan(&path.URL, &path.Path)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		m[path.Path] = path.URL
	}
	defer db.Close()

	return m, nil
}
