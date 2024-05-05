# Devtools (Developer Tools)

* Devtools (Developer Tools) — a set of tools that automate and ease the developer's workflow.
* Additional env variables must be specified to ensure a proper work of devtools (see [.env.template](../../../.env.template)).

## Dashboard backups

* Grafana dashboards can be saved and loaded using 2 commands located in the [dashboards](../../../devtools/dashboards/) directory.
* `GRAFANA_HOST`, `GRAFANA_PORT` and `GRAFANA_TOKEN` environment variables should be exported in advance. `GRAFANA_TOKEN` is Grafana service token (see [official documentaiton](https://grafana.com/docs/grafana/latest/administration/service-accounts/) for more details).
* [Save](../../../devtools/dashboards/save) command saves dashboards from Grafana instance into [grafana/dashboards/](../../../grafana/dashboards/) directory.
* [Load](../../../devtools/dashboards/load) command loads dashboards from [grafana/dashboards/](../../../grafana/dashboards/) directory into Grafana instance.

## Profiling with perf

* Perf — "a tool for using the performance counters subsystem in Linux, and has had various enhancements to add tracing capabilities" [[source](https://perf.wiki.kernel.org/index.php/Main_Page)].
* `pprof` middleware and handlers should be enabled in server configuration. Moreover, its `path_prefix` should be equal "/debug/pprof".

  (see [Reactions Storage - Configuration](../sections/reactions_storage.md#configuration) for more details).

* Run the following command to profile your Reactions Storage instance and generate a Flame Graph:

  ```bash
  ./devtools/perf/fg
  # or with number of seconds to profile (default: 30)
  ./devtools/perf/fg 60
  ```
