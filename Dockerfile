FROM golang:1.25.7-alpine
RUN apk add --no-cache git curl build-base
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", ".air.toml"]