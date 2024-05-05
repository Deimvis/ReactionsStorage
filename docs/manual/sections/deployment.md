# Deployment

* Deployment setup allows you to deploy, start and stop any of supported resources (such as server or database) with a single command in any available environment (local, docker-compose and remotely).

## General Information

### Resources

* All Reactions Storage components are grouped into 4 `resources`: Reactions Storage (server), Database, Monitoring and Simulation.
* Resources are referred with the following keywords
  * `rs` — Reactions Storage (server)
  * `db` — Database (PostgreSQL)
  * `monitoring` — Monitoring (Prometheus, Grafana, Pushgateway)
  * `sim` — Simulation

### Deployment Types

* Each of resources can be independently deployed with any of available `deployment types`. Currently, 3 deployment types are supported: `local`, `docker-compose` and `vm` (virtual machine).
* Each `deployment type` requires a corresponding env file ".env.<deployment_type>". Thus, ".env.local" is used for `local` deployment type and ".env.docker-compose" is used for `docker-compose` deployment type.
* Each `deployment type` corresponds to a single folder within [deploy](../../../deploy/) directory.

### Quick Start

* To start using a particular `deployment type`, one must create a corresponding env file ".env.<deployment_type>" (e.g. ".env.local") located in the root folder (see [.env.template](../../../.env.template)). Afterwards, `resources` can be managed with `deploy/cmd` command.
* `deploy/cmd` command allows to deploy, start and stop different `resources`. All the commands follow this format: `deploy/cmd <deployment_Type> <resource> <action>` (e.g. `deploy/cmd local rs deploy` and `deploy/cmd local rs stop`).

## Available Commands

* All commands are just executables located in the `cmd` directory of the corresponding `deployment_type` (e.g. [deploy/local/cmd](../../../deploy/local/cmd) directory)
* Commands have the following naming format: "\<resource\>_\<action\>" (e.g. ["rs_deploy"](../../../deploy/local/cmd/rs_deploy))
* In order to run the particular command one should run either executable directly or via `deploy/cmd <deployment_type> <resource> <action>`

## Tuning the configuration

* Folder organization differs among `deployment types`.
* E.g. `docker-compose` deployment type has [grafana](../../../deploy/docker-compose/grafana) folder with [datasource.yml](../../../deploy/docker-compose/grafana/datasource.yml) configuration file. Modifying this file will affect future use of this deployment type.
* There is no common rule to configure a particular resource. One should look at the source code of commands and find out how the configuration file is used and what changes he needs to apply. However, in the most of time the answer is obvious. Thus, "grafana" folder contains configuration files for Grafana and "nginx" folder — configuration for Nginx.
