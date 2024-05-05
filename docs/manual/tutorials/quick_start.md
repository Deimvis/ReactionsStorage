# Quick Start

* Install dependencies

```bash
go get
```

* Set environment variables

```bash
cp .env.template .env
# fill .env file
# REMINDER: Do not forget to run PostgreSQL instance at the address specified using the environment variables
```

* Launch server

```bash
. devtools/exenv .env  # export env variables
make run
```

* Run tests

```bash
. devtools/exenv .env # export env variables
make test
```
