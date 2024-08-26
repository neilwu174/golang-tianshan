package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/neilwu174/calculator/internal/handler"
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

func initStatic(mux *http.ServeMux) {
	dir := http.Dir("./")
	fs := http.FileServer(dir)
	mux.Handle("/static/", fs)
}

// func initRoot(mux *http.ServeMux) {
// 	dir := http.Dir("/Users")
// 	fs := http.FileServer(dir)
// 	mux.Handle("/fileSystem/", http.StripPrefix("/fileSystem", neuter(fs)))
// }

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println(r.URL.Path)
		// status, err := os.Stat(r.URL.Path)
		// if err != nil {
		// 	panic(err)
		// }
		// if status.IsDir() {
		// 	log.Println(r.URL.Path, " is a Dir")
		// } else {
		// 	log.Println(r.URL.Path, " is a File")
		// }
		// entries, err := os.ReadDir("./")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// for _, e := range entries {
		// 	w.Write([]byte("http://localhost:8080/fileSystem" + r.URL.Path + e.Name()))
		// }
		// next.ServeHTTP(w, r)
		handler.GetFileSystem(getTemplate(html_filesystem), w, r)
	})
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

func main00aa() {
	err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	r := http.NewServeMux()
	log.Println("Creating route...")
	r.HandleFunc("/home", Index)
	r.HandleFunc("/user", User)
	r.HandleFunc("/explorer", Explorer)
	r.HandleFunc("/explorer/env", Env)
	r.HandleFunc("/explorer/files", redirect)
	r.HandleFunc("/filesystem/{?}", FileSystem)
	initStatic(r)
	// initRoot(r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}
