# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app/cmd

RUN go build -o frontend .

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/cmd/frontend /app

# COPY --from=builder /app/.env /app

CMD [ "/app/frontend" ]