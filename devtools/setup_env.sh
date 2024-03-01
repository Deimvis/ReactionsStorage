#!/usr/bin/env bash

set -e

function sql() {
    psql -U dbrusenin -d reactions -f $1
    return $?
}

sql devtools/sql/destroy_db.sql
sql sql/init_configuration_storage.sql
sql sql/init_reactions_storage.sql
sql devtools/sql/setup_db.sql
