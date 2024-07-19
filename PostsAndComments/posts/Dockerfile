FROM golang:1.17

WORKDIR /go/src/app

COPY . .

RUN go build -o myapp

EXPOSE 8080

CMD ["./myapp"]