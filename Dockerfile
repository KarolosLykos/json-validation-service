FROM golang:1.19.3-alpine AS builder

# Create and change to the app directory.
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN  go build -v -o main cmd/main.go

## Deploy
FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]

EXPOSE 8082
