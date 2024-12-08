#!/usr/bin/env bash

echo "==> Checking that code complies with goimports requirements..."

goimports_files=$(goimports -d -l ./...)
if [[ -n ${goimports_files} ]]; then
    echo 'goimports needs to be run on the following files:'
    echo " ===== "
    echo "${goimports_files}"
    echo " ===== "
    echo "You can use the command: \`make fiximports\` to resolve."
    exit 1
fi
