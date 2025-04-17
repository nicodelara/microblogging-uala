FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN ls -la /app
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tweets ./cmd/tweets

RUN chmod +x tweets

FROM scratch
COPY --from=builder /app/tweets /tweets
EXPOSE 8081
ENTRYPOINT ["/tweets"]
