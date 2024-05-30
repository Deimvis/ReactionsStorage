# How to export Reactions Storage dashboards to your Grafana instance

1. Install GDG tool: https://software.es.net/gdg/docs/gdg/installation/
2. Make sure that the following environment variables are exported: `GRAFANA_HOST`, `GRAFANA_PORT` and `GRAFANA_TOKEN`. `GRAFANA_TOKEN` is Grafana service token (see [official documentaiton](https://grafana.com/docs/grafana/latest/administration/service-accounts/) for more details).
3. Run

   ```bash
   ./devtools/dashboards/load
   ```
* See [Devtools - Dashboard backups](../sections/devtools.md#dashboard-backups) for more details.
