build:
	@go build -o bin/reactions_storage

run: build
	@./bin/reactions_storage $(ARGS)

test:
	@devtools/kportpid 8080 && gotestsum --format dots

# create .env.local before running this target (see deploy/README.md for more details)
test-local:
	@. devtools/exenv .env.local && devtools/kportpid 8080 && gotestsum --format dots
