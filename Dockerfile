FROM golang:alpine

WORKDIR /app
COPY . /app

RUN go build -o main .

EXPOSE 8000

ENTRYPOINT ["./main"]