FROM golang:1.25.1-alpine3.22 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/app/auth-service ./cmd/app

FROM alpine:3.22 AS app
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
COPY --from=builder /app/bin/app/auth-service /usr/local/bin/
WORKDIR /usr/local/bin/
CMD ["./auth-service"]
