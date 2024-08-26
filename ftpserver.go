package main

import (
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

func main_001() {
	factory := &filedriver.FileDriverFactory{
		RootPath: "/Users/xiaoliwu/workspace",
		Perm:     server.NewSimplePerm("root", "root"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     2001,
		Hostname: "::",
	}
	server := server.NewServer(opts)
	server.ListenAndServe()

}
