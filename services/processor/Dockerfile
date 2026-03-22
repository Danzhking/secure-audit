FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /processor ./cmd/processor/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /processor /processor
EXPOSE 9091
CMD ["/processor"]
