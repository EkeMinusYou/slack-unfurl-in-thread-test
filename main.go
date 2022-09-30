package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	SlackToken        string
	VerificationToken string
)

func init() {
	SlackToken = os.Getenv("SLACK_TOKEN")
	VerificationToken = os.Getenv("VERIFICATION_TOKEN")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
