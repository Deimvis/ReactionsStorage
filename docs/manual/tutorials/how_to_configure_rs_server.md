# How to configure a Reactions Storage server

* Reactions Storage configuration includes 2 things: environment variables and server configuration

## 1. Environment variables

* Environment variables are defined in the env files (see [.env.template](../../../.env.template))
* Each env file corresponds to a single `deployment type` (e.g. ".env.local" for `local` deployment type and ".env.docker-compose" for `docker-compose` deployment type). See [Deployment Types](../sections/deployment.md#deployment-types) for more details.

## 2. Server Configuration

* Server configuration includes a single yaml file.
* Configuration example: [server.yaml](../../../configs/server.yaml).
* Configuration file path should be given to the server binary via `--config` flag. Default value: "configs/server.yaml".
* See [Reactions Storage - Configuration](../sections/reactions_storage.md#configuration) for more details.

## Pratical example for running Reactions Storage server locally

1. Create ".env.local" file in the root folder (see [.env.template](../../../.env.template)).

    ```yaml
    DEBUG=1

    PORT=8080
    DATABASE_URL=postgres://dbrusenin@localhost:5432/reactions
    SQL_SCRIPTS_DIR=/Users/dbrusenin/Projects/code/ReactionsStorage/sql/
    ```

2. Create "server.yaml" file at "configs/server.yaml" or update the [example]((../../../configs/server.yaml)).

   ```yaml
   gin:
     general:
       mode: release
       trusted_proxies: []
     middlewares:
       logger:
         enabled: true
       recovery:
         enabled: true
       prometheus:
         enabled: true
         metrics_path: /metrics
       pprof:
         enabled: false
    handlers:
      debug:
        pprof:
          enabled: false
          path_prefix: /debug/pprof
        mem_usage:
          enabled: false
          path: /debug/sys/mem

    pg:
      pool:
        min_conns: 0
        max_conns: 10
        max_conn_lifetime_jitter_s: 1

   ```

3. Launch the server

   `deploy/cmd local rs deploy`
