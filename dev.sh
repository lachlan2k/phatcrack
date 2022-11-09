#!/bin/bash

docker-compose down --remove-orphans
docker-compose -f docker-compose.dev.yml --env-file .dev.env up --build --force-recreate