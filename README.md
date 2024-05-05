# Reactions Storage (Strategies Branch)

**This is a Strategies Branch. There are `.impl_<strategy_name`> directories that represent different implementations strategies: `mvp`, `rt_join`, `no_join`, `async_join`. See [paper](#paper) fore more details.**

_HSE University Diploma Project_

<!-- <div align="center"> -->
A general-purpose service for storing, managing and retrieving user reactions in efficient yet flexible manner.
<!-- </div> -->

## Paper

* TODO: link

## Features

* TODO: list

## Performance

* TODO: Grafana screenshot

## Quick Start

* Install dependencies

```bash
go get
```

* Set environment variables

```bash
cp .env.template .env
# fill .env file
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

## Manual

* Detailed documentation can be found in the [Manual](docs/manual/README.md).

### Reactions Storage Server

* Reactions Storage server — an HTTP web server that handles all incoming requests

* It requires configuration file and environment variables

  (see [Reactions Storage Server](docs/manual/sections/reactions_storage.md) section of [Manual](docs/manual/README.md))

### Deployment (enables running locally, with docker-compose or via remote VMs)

* Deployment setup allows you to deploy, start and stop any of supported resources (such as server or database) with a single command in any available environment (local, docker-compose and remotely)

  (see [Deployment](docs/manual/sections/deployment.md) section of [Manual](docs/manual/README.md))

### Simulation

* Simulation is a tool written in Go for conducting load tests with user-defined workload.

  (see [Simulation](docs/manual/sections/simulation.md) section of [Manual](docs/manual/README.md))
