# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app/cmd

RUN go build -o backend .

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/cmd/backend /app

# COPY --from=builder /app/.env /app

CMD [ "/app/backend" ]