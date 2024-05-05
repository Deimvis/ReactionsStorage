# How to run a fully-fledged service (with database and monitoring subsystems) using Docker

1. Create ".env.docker-compose" file in the root folder (see [.env.template](../../../.env.template)). For example:

    ```yaml
    DEBUG=1

    DATABASE_URL=postgres://dbrusenin:777@db:5432/reactions

    POSTGRES_USER=dbrusenin
    POSTGRES_PASSWORD=777
    POSTGRES_DB=reactions
    ```

2. Apply env variables and create a docker-compose file:

    ```bash
    . devtools/exenv .env.docker-compose
    envsubst < docker-compose.yaml.template > docker-compose.yaml
    ```

3. Tune configurations in [docker-comopse deployment folder](../../../deploy/docker-compose/).
4. Deploy necessary resources. For example:

   ```yaml
    deploy/cmd docker-compose db deploy
    deploy/cmd docker-compose rs deploy
    deploy/cmd docker-compose monitoring deploy
    deploy/cmd docker-compose sim deploy  # actually starts the simulation!
   ```
