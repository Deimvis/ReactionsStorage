# How to run a fully-fledged service (with database and monitoring subsystems) using remote VMs

1. Create ".env.vm" file in the root folder (see [.env.template](../../../.env.template)). For example:

    ```yaml
    DEBUG=1

    PORT=8080
    DATABASE_URL=postgres://dbrusenin:777@123.000.111.222:5433/reactions

    POSTGRES_USER=dbrusenin
    POSTGRES_PASSWORD=777
    POSTGRES_DB=reactions
    
    RS_SSH_HOST=111.111.0.1
    RS_SSH_USER=dbrusenin
    DB_SSH_HOST=222.222.0.2
    DB_SSH_USER=dbrusenin
    MONITORING_SSH_HOST=333.333.0.3
    MONITORING_SSH_USER=dbrusenin
    SIM_SSH_HOST=444.444.0.4
    SIM_SSH_USER=dbrusenin
    ```

2. Apply env variables and create a conf.json file:

    ```bash
    . devtools/exenv .env.vm
    envsubst < deploy/vm/conf.json.template > deploy/vm/conf.json
    ```

3. Add RSA keys required for using SSH on remote VMs to the authentication agent:

    ```bash
    # ssh-add PATH_TO_YOUR_PRIVATE_KEY
    ssh-add ~/.ssh/reactions_storage_rsa
    ```

4. Tune configurations in [vm deployment folder](../../../deploy/vm/).

5. Deploy necessary resources. For example:

   ```yaml
    deploy/cmd vm db deploy
    deploy/cmd vm rs deploy
    deploy/cmd vm monitoring deploy
    deploy/cmd vm sim deploy
    # or all at once
    deploy/cmd vm all deploy
   ```

6. Manage your resources:

   ```yaml
   deploy/cmd vm sim run  # start simulation
   deploy/cmd vm sim stop # stop simulation
   deploy/cmd vm db  ssh   # ssh to database VM
   ```
