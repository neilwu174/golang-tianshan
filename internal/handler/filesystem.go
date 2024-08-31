package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileSystemData struct {
	Directory bool
	Parent    bool
	Link      string
	Name      string
	Img       string
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
		status, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		if status.IsDir() {
			log.Println(path, " is a Dir")
			if err := tmplate.Execute(w, getDirectory(path)); err != nil {
				log.Println(err)
			}
		} else {
			log.Println(path, " is a File")
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(fileBytes)
		}
	}
}

func getDirectory(directory string) []FileSystemData {
	var parentDir = filepath.Dir(directory)
	log.Println("parentDir=", parentDir)
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]FileSystemData, len(entries)+1)
	items[0] = FileSystemData{Directory: true, Parent: true, Link: "/filesystem" + parentDir, Name: parentDir}
	for i, s := range entries {
		log.Println("entry=", i, s.Name())
		items[i+1].Directory = s.IsDir()
		items[i+1].Link = "/filesystem" + directory + "/" + s.Name()
		items[i+1].Name = s.Name()
		items[i+1].Parent = false
		if s.IsDir() {
			items[i+1].Img = "folder.gif"
		} else {
			if filepath.Ext(s.Name()) == ".pdf" {
				items[i+1].Img = "pdf.gif"
			} else if filepath.Ext(s.Name()) == ".zip" {
				items[i+1].Img = "compressed.gif"
			} else if filepath.Ext(s.Name()) == ".mp4" {
				items[i+1].Img = "movie.gif"
			} else {
				items[i+1].Img = "text.gif"
			}
		}
		fmt.Printf("%+v\n", items[i+1])
	}
	// items = append(items, FileSystemData{Directory: true, Link: "/filesystem" + parentDir, Name: parentDir})
	return items
}
