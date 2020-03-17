FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main cmd/corona-exporter/main.go

CMD ["/app/main"]
