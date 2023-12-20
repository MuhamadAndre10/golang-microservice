FROM golang:1.18-alpine as builder
LABEL authors="MUHAMAD ANDRE"

RUN mkdir "app"

COPY . /app

WORKDIR /app

RUN go build -o brokerApp ./cmd

RUN chmod +x ./brokerApp


FROM alpine:latest

RUN mkdir "app"

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]