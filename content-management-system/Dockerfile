FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN go build -o content-management-system .

EXPOSE 8002

CMD ["./content-management-system"]
