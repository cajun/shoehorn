package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/config"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var (
	serverOn = false
	root     = "."
)

func init() {
	flag.BoolVar(&serverOn, "server", false, "set true to run server")
	flag.StringVar(&root, "root", ".", "which dir at the apps located")
}

func On() bool {
	return serverOn
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index handler")
	renderTemplateString(w, index)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("css handler")
	renderTemplateString(w, css)
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("js handler")
	renderTemplateString(w, js)
}

func appHandler(w http.ResponseWriter, r *http.Request, app string) {
	path := fmt.Sprintf("%s/%s", root, app)
	fmt.Printf("app handler: %s", path)
	os.Chdir(path)
	config.LoadConfigs()

	b, err := json.Marshal(config.Processes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplateString(w, string(b))

}
func listHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list handler")
	// find apps and list
	files, err := ioutil.ReadDir(root)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dirs := []map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			path := fmt.Sprintf("%s/%s/shoehorn.cfg", root, file.Name())
			fmt.Println(path)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				dirs = append(dirs, map[string]string{"name": file.Name()})
			}
		}
	}

	b, err := json.Marshal(dirs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderTemplateString(w, string(b))
}

func renderTemplateString(w http.ResponseWriter, tmpl string) {
	_, err := fmt.Fprintf(w, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const lenPath = len("/apps/json/")

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}
func Up() {
	http.HandleFunc("/apps/json/", makeHandler(appHandler))
	http.HandleFunc("/list/json", listHandler)
	http.HandleFunc("/css/application.css", cssHandler)
	http.HandleFunc("/js/application.js", jsHandler)
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Server up on port 9369\n")
	fmt.Printf("Root: %s\n", root)
	http.ListenAndServe(":9369", nil)
}
