# Simulation

* Simulation — a tool written in Go for conducting load tests with user-defined workload.
* It's located in [tests/simulation](../../../tests/simulation/) folder.

## Usage

* `cd tests/simulation`
* `go run main.go --config configs/example.yaml`

## Configuration

* Simulation configuration file allows to describe the workload in detail.
* Configuration example: [example.yaml](../../../tests/simulation/configs/example.yaml).

### Rules (`.rules`)

#### Entities and Topics (`.rules.topics`)

* In essence, the simulation emulates real users and their interactions with entities.
* Each entity belongs to a single group called a `topic`, which could represent a chat in a messenger or a feed in social media.
* `.rules.topics.count` — number of topics
* `.rules.topics.namespace_id` — namespace id for each entity within topic
* `.rules.topics.size` — number of entities within topic
* `.rules.topics.shuffle_per_user` — whether to shuffle topic for each user individually

#### Users (`.rules.users`)

* The simulation begins with users selecting a random topic and initiating their actions on topic entities.
* There are five actions available to users in the simulation: switch topic, scroll, add reaction, remove reaction, and quit.
* `.rules.users.count` — number of users
* `.rules.users.action_probs` — user action probabilities (essentially, they are just weights / their sum can be $\gt 100$)
* `.rules.users.screen.visible_entities` — number of visible entities visible on a screen (can also simulate pagination in entity loading).
* `.rules.users.app.background.refresh_reactions.timer_in_turns` — number of turns used as a timer to refresh visible entities (to handle cases where users do not scroll and reactions data for visible entities should be updated).
* `rules.users.turn_start_skew_ms` — duration in milliseconds for users to choose the time to start new turn (each user will start turn at random point of time inside interval \[turn_start_t, turn_start_t + turn_start_skew_ms)).
* `rules.users.id_template` — template string for generation of user id (useful for ensuring that different simulation launches do or do not interfere with each other). Must contain `%d` for counter.

#### Turns (`.rules.turns`)

* Simulation is organized into multiple turns, each representing a unit of time.
* During a single turn, a client performs exactly one action and then waits for the end of the current turn if it is not yet complete. Each turn has a minimum duration specified in milliseconds and ends when both the user completed their action and the minimum duration has passed.
* `.rules.turns.count` — total number of turns for simulation to run
* `.rules.turns.min_duration_ms` — minimum duration for each turn in milliseconds
