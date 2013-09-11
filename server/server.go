package server

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/command"
	"github.com/cajun/shoehorn/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
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
	renderTemplateString(w, index)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplateString(w, css)
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplateString(w, js)
}

func appHandler(w http.ResponseWriter, r *http.Request, app string) {
	path := fmt.Sprintf("%s/%s", root, app)
	os.Chdir(path)
	config.LoadConfigs()

	b, err := json.Marshal(config.Processes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplateString(w, string(b))

}

func listHandler(w http.ResponseWriter, r *http.Request) {
	// find apps and list
	files, err := ioutil.ReadDir(root)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dirs := []map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			path := fmt.Sprintf("%s/%s/shoehorn.cfg", root, file.Name())
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

func commandHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Command Handler")
	base := r.URL.Path[len("/commands/"):]
	opts := strings.Split(base, "/")
	fmt.Printf("Commands: %v", opts)
	site := opts[0]
	process := opts[1]
	cmd := opts[2]
	path := fmt.Sprintf("%s/%s", root, site)

	os.Chdir(path)
	command.MkDirs()
	config.LoadConfigs()

	old := os.Stdout
	re, wr, _ := os.Pipe()
	os.Stdout = wr

	command.ParseCommand([]string{process, cmd})

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, re)
		outC <- buf.String()
	}()

	wr.Close()
	os.Stdout = old
	out := <-outC

	clean := strings.Replace(out, "\n", "</br>", -1)
	o := map[string]string{"status": "ok", "output": clean}

	b, err := json.Marshal(o)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplateString(w, string(b))
}

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.FormValue("repo")
	os.Chdir(root)
	opts := []string{"clone", repo}

	cmd := exec.Command("git", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	http.Redirect(w, r, "/", http.StatusFound)
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
	http.HandleFunc("/clone", cloneHandler)
	http.HandleFunc("/commands/", commandHandler)
	http.HandleFunc("/apps/json/", makeHandler(appHandler))
	http.HandleFunc("/list/json", listHandler)
	http.HandleFunc("/css/application.css", cssHandler)
	http.HandleFunc("/js/application.js", jsHandler)
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Server up on port 9369\n")
	fmt.Printf("Root: %s\n", root)
	http.ListenAndServe(":9369", nil)
}
