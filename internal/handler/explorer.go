package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func Navinate(tmplate *template.Template, w http.ResponseWriter, r *http.Request) {
	data := getEnv()
	if err := tmplate.Execute(w, data); err != nil {
		log.Println(err)
	}
	getEnv()
}
func GetEnv(tmplate *template.Template, w http.ResponseWriter, r *http.Request) {
	data := getEnv()
	if err := tmplate.Execute(w, data); err != nil {
		log.Println(err)
	}
	getEnv()
}

func getEnv() *map[string]interface{} {
	data := make(map[string]interface{})
	envs := os.Environ()
	for _, e := range envs {
		ar := strings.Split(e, "=")
		data[ar[0]] = ar[1]
		// log.Println("value:", data)
	}
	return &data
}
