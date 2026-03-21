FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /api
EXPOSE 8081
CMD ["/api"]
