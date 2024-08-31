package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/neilwu174/calculator/internal/handler"

	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	layoutsDir      = "templates/layouts"
	templatesDir    = "templates"
	extension       = "/*.html"
	index           = "index.html"
	user            = "user.html"
	html_explorer   = "explorer.html"
	html_env        = "env.html"
	html_filesystem = "filesystem.html"
)

var (
	//go:embed templates/* templates/layouts/*
	files     embed.FS
	templates map[string]*template.Template
)

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}
	for _, tmpl := range tmplFiles {
		log.Printf("Scan %s...\n", tmpl.Name())
		if tmpl.IsDir() {
			log.Printf("%s is a Dir\n", tmpl.Name())
			continue
		}

		log.Printf("%s is a File, Parse %s\n", tmpl.Name(), tmpl.Name())
		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

func getTemplate(name string) *template.Template {
	t, ok := templates[name]
	if !ok {
		log.Printf("template %s not found", name)
		panic(ok)
	}
	return t
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("loading index.html")
	handler.LoadIndex(getTemplate(index), w, r)
}
func User(w http.ResponseWriter, r *http.Request) {
	log.Println("loading user")
	handler.ApplyUser(getTemplate(user), w, r)
}

func Explorer(w http.ResponseWriter, r *http.Request) {
	log.Println("loading explorer")
	handler.Navinate(getTemplate(html_explorer), w, r)
}
func Env(w http.ResponseWriter, r *http.Request) {
	log.Println("loading explorer")
	handler.GetEnv(getTemplate(html_env), w, r)
}
func FileSystem(w http.ResponseWriter, r *http.Request) {
	log.Println("loading filesystem")
	handler.GetFileSystem(getTemplate(html_filesystem), w, r)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	log.Println("Redirect ...")
	http.Redirect(w, r, "/", 301)
}

var routes = []route{
	newRoute("GET", "/", home),
	newRoute("GET", "/xyz/(.*)", sink),
	newRoute("GET", "/groups/([^/]+)/people", peopleInGroupHandler),
	newRoute("GET", "/filesystem/(.*)", FileSystem),
	newRoute("GET", "/filesystem", FileSystem),
	newRoute("GET", "/home", Index),
	newRoute("GET", "/explorer", Explorer),
	newRoute("GET", "/explorer/env", Env),
}

func getRouts() []route {
	dir := http.Dir("./")
	fs := http.FileServer(dir)
	var routes = append(routes, newRoute("GET", "/static/(.*)", fs.ServeHTTP))
	return routes
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range getRouts() {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		fmt.Println("matches=", matches)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			fmt.Println("route=", route.handler)
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HOME\n")
}
func peopleInGroupHandler(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "Group handler: %s\n", slug)
}

func fileSystemHandler(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "File System handler: %s\n", slug)
}

func sink(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "Sink %s\n", slug)
}

func main() {
	LoadTemplates()
	http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(Serve))
}
