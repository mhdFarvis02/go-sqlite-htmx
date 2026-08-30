FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/app .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/app /app/app

EXPOSE 8080

CMD ["/app/app"]
