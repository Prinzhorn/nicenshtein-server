FROM golang:alpine

RUN apk add --no-cache git
RUN go get "github.com/prinzhorn/nicenshtein-server"
RUN apk del git

WORKDIR "/go/src/github.com/prinzhorn/nicenshtein-server"
CMD ["/go/bin/nicenshtein-server"]
EXPOSE 8080