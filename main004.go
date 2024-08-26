package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	layoutsDir01   = "templates/layouts"
	templatesDir01 = "templates"
	extension01    = "/*.html"
)

var (
	//go:embed templates/* templates/layouts/*
	files01     embed.FS
	templates01 map[string]*template.Template
)

var routes = []route{
	newRoute("GET", "/", home),
	newRoute("GET", "/(.*)", sink),
	newRoute("GET", "/groups/([^/]+)/people", peopleInGroupHandler),
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
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
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

func sink(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "Sink %s\n", slug)
}

func main() {
	LoadTemplates01()
	http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(Serve))
}

func LoadTemplates01() error {
	if templates01 == nil {
		templates01 = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files01, templatesDir01)
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
		pt, err := template.ParseFS(files01, templatesDir01+"/"+tmpl.Name(), layoutsDir01+extension01)
		if err != nil {
			return err
		}

		templates01[tmpl.Name()] = pt
	}
	return nil
}
