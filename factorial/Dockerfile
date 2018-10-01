FROM golang:1.9

# setup go build enviroment
RUN go get -u github.com/kardianos/govendor
RUN mkdir -p /go/src/github.com/pintjuk/faas/factorial
COPY . /go/src/github.com/pintjuk/faas/factorial
EXPOSE 8080

WORKDIR /go/src/github.com/pintjuk/faas/factorial
RUN govendor fetch +m
WORKDIR /go/src/github.com/pintjuk/faas/factorial/cmd/factoriald
RUN go build
RUN ls

CMD ["./factoriald"]