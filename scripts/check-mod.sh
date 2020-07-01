#!/usr/bin/env bash

echo "==> Checking that the module is clean..."

go mod download
tidy=$(go mod tidy -v 2>&1)
if [[ ${tidy} ]]; then
    echo 'Extranenous dependencies need removed.'
    echo " ===== "
    echo "${tidy}"
    echo " ===== "
    echo "You can use the command: \`make tidy\` to remove these dependencies."
    git checkout go.*
    exit 1
fi
