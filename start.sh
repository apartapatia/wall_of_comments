#!/bin/bash

DB_TYPE="$1"

if [ "$DB_TYPE" = "redis" ]; then
    docker-compose up -d app redis
elif [ "$DB_TYPE" = "postgres" ]; then
    docker-compose up -d app postgres
else
    echo "Unknown database type: $DB_TYPE"
fi
