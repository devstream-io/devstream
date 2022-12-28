FROM golang:alpine AS build-env
WORKDIR $GOPATH/src/github.com/[[.Repo.Owner]]/[[.Repo.Name]]
COPY . .
RUN apk add git
RUN go get ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/[[.AppName]]/main.go

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/[[.Repo.Owner]]/[[.Repo.Name]]/app /app/
CMD ["./app"]
USER 1000
EXPOSE 8080/tcp
