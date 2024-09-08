#!/bin/bash

cloc ./agent ./api ./common ./frontend --exclude_dir=node_modules --not-match-f=types.ts --exclude-ext=json
