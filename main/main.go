package main

import (
	"log"
	"net/http"

	"glavan.tech/privateNotes"
)

func main() {
	http.HandleFunc("/", privateNotes.PrivateNotes)
	log.Fatal(http.ListenAndServe(":80", nil))
}
