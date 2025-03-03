FROM golang:1.24.0 AS builder
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o client ./cmd/client

FROM scratch

COPY --from=builder /usr/src/app/client /client

ENTRYPOINT ["/client"]
