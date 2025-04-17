FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o users ./cmd/users

RUN chmod +x users

FROM scratch
COPY --from=builder /app/users /users
EXPOSE 8082
ENTRYPOINT ["/users"]
