package main

import (
	"flag"
	"net/http"
)

// Minimal target web server used as the reference application in the benchmark.
// It returns the same static page for every request, so the measured cost is the cost of the agent.
func main() {
	address := flag.String("address", "127.0.0.1:8080", "The address to listen on")
	flag.Parse()

	body := []byte("<html><body><h1>It works!</h1></body></html>")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(body)
	})
	http.ListenAndServe(*address, nil)
}
