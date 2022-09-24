package main

import (
	"log"
	"net/http"

	// Blank-import the function package so the init() runs
	"glavan.tech/privateNotes"
)

func main() {
	http.HandleFunc("/", privateNotes.PrivateNotes)
	log.Fatal(http.ListenAndServe(":80", nil))
}
