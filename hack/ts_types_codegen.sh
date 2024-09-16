#!/bin/bash

# go install github.com/tkrajina/typescriptify-golang-structs/tscriptify@latest
# go get github.com/tkrajina/typescriptify-golang-structs

cd common
tscriptify -package=github.com/lachlan2k/phatcrack/common/pkg/apitypes -target=../frontend/src/api/types.ts -interface ./pkg/apitypes/*

sed -i 's/key: int/key: number/' ../frontend/src/api/types.ts

cd ../frontend
npm run lint