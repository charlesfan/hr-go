FROM golang:1.24.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod download
RUN go build -o /app/cmd/server /app/server/server.go
EXPOSE 8080
ENV GIN_MODE=release
CMD ["/app/cmd/server"]
