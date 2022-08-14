FROM golang:1.19
WORKDIR /go/src/gocc

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY cmd ./cmd
COPY *.go ./

COPY Makefile ./
COPY test.sh ./

RUN make test


