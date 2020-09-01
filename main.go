//
// This program is a prototype to learm more about jir webhooks
//
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"net/smtp"
)

const(
	FROM = "m.darin@email.su"
	PASSWORD = "x_wK,(Tf6ff)2L_"
)
//-------------------------------
// main driver
//-------------------------------
func main() {
	// configure handlers
	http.HandleFunc("/", handle_hello)
	http.HandleFunc("/hi", handle_hi)
	http.HandleFunc("/webhookreceiver", handle_jira_webhook)
	http.HandleFunc("/world", handle_world)

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

func handle_world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[world] webhook machined!")
}

func handle_jira_webhook(w http.ResponseWriter, r *http.Request) {
	// yandex
	// smtp.yandex.ru:465
	// http://ilyakhasanov.ru/baza-znanij/prochee/nuzhno-znat/139-nastrojki-otpravki-pochty-cherez-smtp
	// how to resolve EOF error
	// https://stackoverflow.com/questions/11662017/smtp-sendmail-fails-after-10-minutes-with-eof

	log.Println("auth")
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", FROM, PASSWORD, "smtp.yandex.ru")

	log.Println("prepare and send")
	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"m.darin.comco@yandex.ru"}
	msg := []byte("To: m.darin.comco@yandex.ru\r\n" +
		"Subject: JIRA webhook is machined\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err := smtp.SendMail("smtp.yandex.ru:587", auth, FROM, to, msg)
	if err != nil {
		log.Print("Error: ")
		log.Fatal(err)
	} else {
		log.Println("Successfull done!")
	}
	fmt.Fprintf(w, "webhook is machined!")
}
