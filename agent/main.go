package main

import "github.com/lucacoratu/disertatie/agent/server"

func main() {
	proxyServer := server.AgentServer{}
	err := proxyServer.Init()
	if err != nil {
		return
	}
	proxyServer.Run()
}
