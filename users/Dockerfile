FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o users .

EXPOSE 8000

CMD ["./users"]
