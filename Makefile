build:
	@go build -o bin/reactions_storage

run: build
	@./bin/reactions_storage $(ARGS)

test:
	@devtools/kportpid 8080 && gotestsum --format dots
