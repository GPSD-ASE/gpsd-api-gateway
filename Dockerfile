FROM golang:latest

WORKDIR /app

RUN apt-get update && apt-get install -y curl lsof gcc

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o api_gateway ./internal/cmd

EXPOSE 3000

CMD ["./api_gateway"]