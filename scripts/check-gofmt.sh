#!/usr/bin/env bash

echo "==> Checking that code complies with gofmt requirements..."

gofmt_files=$(gofmt -l -s `find . -name '*.go' | grep -v vendor`)
if [[ -n ${gofmt_files} ]]; then
    echo 'gofmt needs to be run on the following files:'
    echo " ===== "
    echo "${gofmt_files}"
    echo " ===== "
    echo "You can use the command: \`make fmt\` to resolve."
    exit 1
fi
