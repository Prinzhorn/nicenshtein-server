FROM golang:alpine
RUN go install github.com/prinzhorn/nicenshtein-server
CMD ["/go/bin/nicenshtein-server"]
EXPOSE 8080