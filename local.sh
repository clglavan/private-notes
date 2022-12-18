#!/bin/sh
GOOGLE_APPLICATION_CREDENTIALS=key.json FUNCTION_TARGET=privateNotes go run cmd/main.go

# this assumes local development simulating cloud functions

# alternative is to build your docker image with correct arguments. 
# for local development please check private_notes.go for the lines with "// for local development" and uncomment them
