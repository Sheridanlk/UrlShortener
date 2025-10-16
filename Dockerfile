FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/url-shortener .
COPY /config/config.yaml ./config.yaml

ENTRYPOINT ["/app/url-shortener"]