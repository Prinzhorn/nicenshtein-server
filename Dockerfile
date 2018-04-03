FROM golang:alpine
ADD . /go/src/github.com/prinzhorn/nicenshtein-server
# go get nicenshtein
RUN go install github.com/prinzhorn/nicenshtein-server
CMD ["/go/bin/nicenshtein-server", "/go/src/github.com/prinzhorn/nicenshtein-server/10_million_password_list_top_100000.txt"]
EXPOSE 8080