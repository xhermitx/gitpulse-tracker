# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o githupApp .

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/githupApp /app

# COPY --from=builder /app/.env /app

CMD [ "/app/githupApp" ]