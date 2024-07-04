# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o logger .

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/logger /app

CMD [ "/app/logger" ]