FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o simple-api main.go

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/simple-api .

EXPOSE 8080
CMD ["./simple-api"]