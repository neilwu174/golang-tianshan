package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type FileSystemData struct {
	Directory bool
	Name      string
}

func GetFileSystem(tmplate *template.Template, w http.ResponseWriter, r *http.Request) {
	log.Println("GetFileSystem-checking", r.URL.Path)
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	if r.URL.Path == "/filesystem" {
		if err := tmplate.Execute(w, getDirectory("/Users")); err != nil {
			log.Println(err)
		}
	} else {
		path := r.URL.Path[len("/filesystem"):len(r.URL.Path)]
		if err := tmplate.Execute(w, getDirectory(path)); err != nil {
			log.Println(err)
		}
	}
}

func getDirectory(directory string) []FileSystemData {
	status, err := os.Stat(directory)
	if err != nil {
		panic(err)
	}
	if status.IsDir() {
		log.Println(directory, " is a Dir")
	} else {
		log.Println(directory, " is a File")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]FileSystemData, len(entries))
	for i, s := range entries {
		log.Println("entry=", i, s.Name())
		items[i].Directory = s.IsDir()
		items[i].Name = "/filesystem" + directory + "/" + s.Name()
	}
	return items
}
