FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/app .

FROM alpine:3.20

RUN apk add --no-cache bash git tmux

RUN mkdir -p /workspace /app

WORKDIR /workspace

COPY --from=builder /out/app /app/app

COPY ./start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start

EXPOSE 8080

# Ensure the data dir exists before the app starts, then launch the tmux wrapper.
RUN mkdir -p /workspace/data

CMD ["/start"]