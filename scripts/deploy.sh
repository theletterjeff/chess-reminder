#!/bin/bash

set -e

BASE_DIR="$(pwd)"
SRC_DIR="$BASE_DIR/src"
CDK_DIR="$BASE_DIR/cdk"

cd "$SRC_DIR"
GOOS=linux GOARCH=amd64 go build -o main .
zip function.zip main

# remove build artifacts on exit
function cleanup {
    rm "$SRC_DIR/main"
    rm "$SRC_DIR/function.zip"
}
trap "cleanup $?" EXIT

cd "$CDK_DIR"

npm install
npm run build

cdk deploy
