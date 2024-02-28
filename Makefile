build:
	@go build -o bin/reactions_storage

run: build
	@./bin/reactions_storage $(ARGS)

test:
	@go test -v ./...
