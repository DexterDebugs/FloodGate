FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download     

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o floodgate ./cmd/floodgate

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/floodgate .

COPY configs/ ./configs/

EXPOSE 8080

CMD ["./floodgate"]