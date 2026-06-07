# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY . .

RUN go build -o /aevra .

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /aevra /aevra
COPY --from=builder /app/templates /app/templates

EXPOSE 8091

USER nonroot:nonroot
ENTRYPOINT ["/aevra"]
