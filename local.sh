#!/bin/sh
GOOGLE_APPLICATION_CREDENTIALS=key.json FUNCTION_TARGET=privateNotes go run cmd/main.go

# for local developer please check private_notes.go for the lines with "// for local development" and uncomment them
