FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /out/api /api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]