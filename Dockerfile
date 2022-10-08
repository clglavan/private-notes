FROM golang:1.16-alpine

ENV ENVIRONMENT="docker"
ENV GCP_PROJECT=""
ENV GCP_REGION=""
ENV GCP_BUCKET_NAME="private-notes"

WORKDIR /private-notes

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /private-notes/

RUN GOOS=linux GARCH=amd64 go build -o /private-notes/privateNotes ./main
RUN chmod 755 /private-notes/privateNotes

EXPOSE 80

CMD ["/private-notes/privateNotes" ]
