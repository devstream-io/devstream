FROM golang:1.17 as builder

WORKDIR /devstream

# cache deps before building so that source changes don't invalidate our downloaded layer
COPY go.mod go.sum /devstream/

RUN go mod download

CMD ["./build_linux_amd64.sh"]
