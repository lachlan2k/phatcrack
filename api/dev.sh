#!/bin/bash

docker-compose down --remove-orphans
docker-compose -f docker-compose.dev.yml up --build --force-recreate