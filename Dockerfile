FROM golang:latest AS builder

WORKDIR /build
ADD go.mod .

COPY . .

RUN go build -o weather cmd/main.go

FROM ubuntu:latest

WORKDIR /build

COPY --from=builder /build/config /build/config
COPY --from=builder /build/weather /build/weather

CMD ["./weather", "IPV4One=1"]
