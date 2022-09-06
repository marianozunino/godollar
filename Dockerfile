# Build Stage
FROM golang:1.19 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app.bin

# Final Stage
FROM golang:1.19

WORKDIR /app

COPY --from=builder /build/app.bin .
COPY --from=builder /build/*.sql .

ENV GIN_MODE=release
ENV PORT=8080
EXPOSE 8080

CMD ["/app/app.bin"]


