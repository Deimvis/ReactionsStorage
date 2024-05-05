# How to conduct a load test with custom workload

* Load testing — a type of performance testing that simulates a real-world load.
* We developed a [simulation](../../../tests/simulation/) tool for these purposes. See [Simulation](../sections/simulation.md) section for more details.

## Practical example to launch simulation using remote VMs

1. Configure the `vm` deployment type following instructions from [How to run a fully-fledged service (with database and monitoring subsystems) using remote VMs](how_to_run_service_using_remote_vms.md).

2. Create simulation configuration yaml file inside of [simulation](../../../deploy/vm/simulation/) folder. For example, create `load-test.yaml`:
   ```yaml
   seed: 123
   rules:
   turns:
     count: 180
     min_duration_ms: 1000

   users:
     count: 10
     turn_start_skew_ms: 900

     id_template: user_%06d
    
     screen:
       visible_entities_count: 3
    
     app:
       background:
         refresh_reactions:
         timer_in_turns: 3

     action_probs:
        switch_topic: 1
        scroll: 2
        add_reaction: 3
        remove_reaction: 4
        quit: 5

     topics:
       - count: 10
         namespace_id: namespace
         size: 1000
         shuffle_per_user: true

   server:
     host: 111.111.0.1
     port: 8080
     ssl: false

   prometheus_pushgateway:
     host: 222.222.0.2
     port: 9091
     ssl: false
     push_interval_s: 15

   ```

3. Use `deploy/cmd vm sim run --config configs/NAME_OF_YOUR_CONFIG_FILE` to start simulation. For example: `deploy/cmd vm sim run --config configs/load-test.yaml`.
  * You can use `devtools/load-test/launch` to start simulation on remote VMs with automated metrics capturing (result will be stored in metrics.json)
  
    **!!!** Make sure all required devtools env variables are exported before launch (see [.env.template](../../../.env.template)).

    **!!!** Configuration file should be named `load-test.yaml`.
