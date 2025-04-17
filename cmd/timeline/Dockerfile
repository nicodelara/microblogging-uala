# Etapa de construcción
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o timeline ./cmd/timeline

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/timeline /app/timeline

EXPOSE 8083

ENV TIMELINE_PORT=8083

ENTRYPOINT ["/app/timeline"]
