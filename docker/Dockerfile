FROM golang:1.19-alpine

WORKDIR /private-notes

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /private-notes/

RUN go build -o /private-notes/privateNotes ./main
RUN chmod 755 /private-notes/privateNotes

EXPOSE 80

CMD ["/private-notes/privateNotes" ]
