# Reactions Storage Server

* Reactions Storage server — an HTTP web server that handles all incoming requests

## Usage

* `. devtools/exenv .env`
* `go run main.go --config configs/server.yaml`

## Configuration

### General

* Server configuration includes a single yaml file.
* Configuration example: [server.yaml](../../../configs/server.yaml).
* Configuration file path should be given to the server binary via --config flag. Default value: "configs/server.yaml".

### Settings

* All available settings can be found in [server.go](../../../src/configs/server.go).
* `gin` — represents gin related settings (mode, middleware, http handlers, etc)
  * `general`
    * `mode`: debug/release/... — one of gin modes.
    * `trusted_proxies` — trusted proxies setting
  * `middlewares`
    * `logger`
      * `enabled`: true/false
    * `recovery`
      * `enabled`: true/false
    * `prometheus`
      * `enabled`: true/false
      * `metrics_path` — URL path for metrics handler
    * `pprof`
      * `enabled`: true/false
  * `handlers`
    * `debug`
      * `pprof`
        * `enabled`: true/false (.gin.middlewares.pprof.enabled should also be enabled)
        * `path_prefix` — URL path prefix for pprof handlers
      * `mem_usage`
        * `enabled`: true/false
        * `path` — URL path for mem usage handler
* `pg`
  * `pool`
    * `min_conns` — min active connections to PostgreSQL
    * `max_conns` — max active connections to PostgreSQL
    * `max_conns_lifetime_jitter_s` — max random delay before closing active connection
