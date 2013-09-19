package server

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/command"
	"github.com/cajun/shoehorn/config"
	"github.com/cajun/shoehorn/logger"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	serverOn = false
)

func init() {
	flag.BoolVar(&serverOn, "server", false, "set true to run server")
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
	path := fmt.Sprintf("%s/%s", command.Root(), app)
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
	files, err := ioutil.ReadDir(command.Root())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dirs := []map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			path := fmt.Sprintf("%s/%s/shoehorn.cfg", command.Root(), file.Name())
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
	logger.Log(fmt.Sprintln("Command Handler"))
	base := r.URL.Path[len("/commands/"):]
	opts := strings.Split(base, "/")
	logger.Log(fmt.Sprintf("Commands: %v", opts))
	site := opts[0]
	process := opts[1]
	cmd := opts[2]
	path := fmt.Sprintf("%s/%s", command.Root(), site)

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
	command.Install(r.FormValue("repo"))
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

func before(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log(fmt.Sprintln(r.URL.Path))
		fn(w, r)
	}
}

func Up() {
	http.HandleFunc("/clone", before(cloneHandler))
	http.HandleFunc("/commands/", before(commandHandler))
	http.HandleFunc("/apps/json/", before(makeHandler(appHandler)))
	http.HandleFunc("/list/json", before(listHandler))
	http.HandleFunc("/css/application.css", before(cssHandler))
	http.HandleFunc("/js/application.js", before(jsHandler))
	http.HandleFunc("/", before(indexHandler))

	logger.Log(fmt.Sprintf("Server up on port 9369\n"))
	logger.Log(fmt.Sprintf("Root: %s\n", command.Root()))

	http.ListenAndServe(":9369", nil)
}
