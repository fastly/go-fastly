#!/usr/bin/env bash

GO_FOR_BUILD="${1}"

echo "==> Checking that the module is clean..."

tidy=$(${GO_FOR_BUILD} mod tidy -v 2>&1 | grep -v "^go: downloading")
if [[ ${tidy} ]]; then
    echo 'Extranenous dependencies need removed.'
    echo " ===== "
    echo "${tidy}"
    echo " ===== "
    echo "You can use the command: \`make tidy\` to remove these dependencies."
    git checkout go.*
    exit 1
fi
