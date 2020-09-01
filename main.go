//
// This program is a prototype to learm more about jir webhooks
//
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

//-------------------------------
// main driver
//-------------------------------
func main() {
	// configure handlers
	http.HandleFunc("/", handle_hello)
	http.HandleFunc("/hi", handle_hi)
	http.HandleFunc("/webhookreceiver", handle_jira_webhook)

	// start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//-------------------------------
// handlers
//-------------------------------
func handle_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path[1:]))
}

func handle_hi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

func handle_jira_webhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "webhook machined!")
}
