FROM golang:1.23

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./tmp/main -buildvcs=false

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]