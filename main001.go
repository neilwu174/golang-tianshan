package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

var (
	//go:embed resources
	res embed.FS

	pages = map[string]string{
		"/":             "resources/index.html",
		"/user-profile": "resources/userProfile.html",
	}
)

func main001() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started\n")
		page, ok := pages[r.URL.Path]
		if !ok {
			log.Fatal(ok)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("Started-Reading\n")
		tpl, err := template.ParseFS(res, page)
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Started-Sending,%p\n", tpl)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		log.Printf("Started-Sending-001\n")
		if err := tpl.Execute(w, data); err != nil {
			log.Fatal(err)
			return
		}
	})
	http.HandleFunc("/user-profile", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started-user-profile\n")
		page, ok := pages[r.URL.Path]
		if !ok {
			log.Fatal(ok)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("Started-Reading\n")
		tpl, err := template.ParseFS(res, page)
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Started-Sending,%p\n", tpl)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		log.Printf("Started-Sending-001\n")
		if err := tpl.Execute(w, data); err != nil {
			log.Fatal(err)
			return
		}
	})

	http.FileServer(http.FS(res))

	log.Println("server started...")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}
