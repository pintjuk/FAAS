FROM golang:1.9

# go stuff
RUN go get -u github.com/kardianos/govendor && \
go get -u -v github.com/codeskyblue/fswatch
RUN mkdir -p /go/src/github.com/pintjuk/faas/addition
WORKDIR /go/src/github.com/pintjuk/faas/addition

EXPOSE 8080

CMD govendor fetch +m ; go run cmd/*/*.go