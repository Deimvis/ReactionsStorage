# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY sql ./sql
COPY src ./src
COPY main.go ./
COPY Makefile ./

# ENTRYPOINT ["/bin/bash"]
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/reactions_storage

EXPOSE 8080
ENV SQL_SCRIPTS_DIR=/app/sql
CMD ["./bin/reactions_storage"]
