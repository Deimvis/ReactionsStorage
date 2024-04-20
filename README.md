# Reactions Storage

_HSE University Diploma Project_

## Quick Start

* Install dependencies
```bash
go get
```

* Set environment variables

```bash
cp .env.template .env
# fill .env file
. devtools/exenv .env  # export env variables
```

* Launch server

```bash
make run
```

## Development

### Testing

```bash
. devtools/exenv .env
make test
```

### Integration with air

```bash
go install github.com/cosmtrek/air@latest
air
```
