FROM golang:1.16-alpine

WORKDIR /private-notes

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /private-notes/

RUN cd cmd
RUN GOOS=linux GARCH=amd64 go build -o /private-notes/main
RUN chmod 755 /private-notes/main

# RUN ls -la

EXPOSE 80

CMD ["/bin/sh","-c","cd /private-notes; /private-notes/main" ]
