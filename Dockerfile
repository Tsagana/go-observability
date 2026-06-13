FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api && go build -o worker ./cmd/worker && go build -o publisher ./cmd/publisher

EXPOSE 8080
