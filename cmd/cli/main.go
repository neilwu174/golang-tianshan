package main

import (
	"github.com/neilwu174/calculator/internal/random"

	// Import the color package.

	"net/http"

	"github.com/fatih/color"
)

type Todo struct {
	Title string
	Done  bool
}

type HttpInfo struct {
	Host   string
	Method string
}

type TodoPageData struct {
	PageTitle string
	Info      HttpInfo
	Todos     []Todo
}

func main() {
	green := color.New(color.FgGreen)
	green.Printf("Your lucky number is %d!\n", random.Number())

	//serve static files
	fs := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	// tmpl := template.Must(template.ParseFiles("static/html/layout.html"))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	dump, _ := httputil.DumpRequest(r, true)
	// 	green.Printf("%q", dump)
	// 	green.Printf("host=%s", r.Host)
	// 	green.Printf("method=%s", r.Method)

	// 	data := TodoPageData{
	// 		PageTitle: "My TODO list",
	// 		Info:      HttpInfo{Host: r.Host, Method: r.Method},
	// 		Todos: []Todo{
	// 			{Title: "Task 1111000", Done: false},
	// 			{Title: "Task 2222", Done: true},
	// 			{Title: "Task 3222", Done: true},
	// 		},
	// 	}
	// 	tmpl.Execute(w, data)
	// })
	http.ListenAndServe(":80", nil)
}
