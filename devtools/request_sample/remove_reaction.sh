#!/usr/bin/env bash

curl -X DELETE localhost:8080/reactions -d '{"namespace_id": "namespace", "entity_id": "entity", "reaction_id": "reaction", "user_id": "user"}' -v
