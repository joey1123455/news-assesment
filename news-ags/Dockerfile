FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN go build -o news-ags .

EXPOSE 8001

CMD ["./news-ags"]
