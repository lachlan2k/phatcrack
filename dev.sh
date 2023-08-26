#!/bin/bash
rm -r api/tmp
mkdir api/tmp
docker-compose -f docker-compose.dev.yml down --remove-orphans
docker-compose -f docker-compose.dev.yml --env-file .dev.env up --build --force-recreate