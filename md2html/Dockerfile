FROM golang:1.9

# go stuff
RUN go get -u github.com/kardianos/govendor
RUN mkdir -p /go/src/github.com/pintjuk/faas/md2html
WORKDIR /go/src/github.com/pintjuk/faas/md2html

EXPOSE 8080

CMD govendor fetch +m ; go run cmd/*/*.go