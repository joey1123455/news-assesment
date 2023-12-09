FROM golang:1.21.3

WORKDIR /go/src/app
COPY . .

RUN go clean -modcache

RUN go mod tidy
RUN go build -o users .

EXPOSE 8000

CMD ["./users"]
