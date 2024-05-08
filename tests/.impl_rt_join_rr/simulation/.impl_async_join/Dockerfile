# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src
COPY tests/simulation ./tests/simulation

WORKDIR /app/tests/simulation
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/simulation

ENTRYPOINT ["./bin/simulation"]
