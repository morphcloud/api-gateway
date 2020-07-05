FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o app .

FROM scratch
COPY --from=builder /build/app /app
CMD ["/app"]
