# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o profilerApp .

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/profilerApp /app

# COPY --from=builder /app/.env /app

CMD [ "/app/profilerApp" ]