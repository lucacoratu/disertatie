package main

import (
	"github.com/lucacoratu/disertatie/api/server"
)

func main() {
	server := server.APIServer{}
	err := server.Init()
	if err != nil {
		return
	}
	server.Run()
}
