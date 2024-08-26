package handler

import (
	"html/template"
	"log"
	"net/http"
)

func ApplyUserProfile(tmplate *template.Template, w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["Name"] = "John Doe Profile"
	data["Email"] = "johndoe@email.com"
	data["Address"] = "Fake Street, 123"
	data["PhoneNumber"] = "654123987"

	if err := tmplate.Execute(w, data); err != nil {
		log.Println(err)
	}
}
