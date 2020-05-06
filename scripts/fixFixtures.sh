#!/bin/bash
set -e

FASTLY_TEST_SERVICE_ID=$1
DEFAULT_TEST_SERVICE_ID="7i6HN3TK9wS159v2gPAZ8A"
FIXTURESDIR="$(pwd)/fastly/fixtures/"

if [[ -z ${FASTLY_TEST_SERVICE_ID} && -z $1 ]]; then
  echo "You must supply a service ID as the first argument"
  exit
fi

for file in $(grep --recursive --files-with-matches "${FASTLY_TEST_SERVICE_ID}" "${FIXTURESDIR}")
do
  sed -i "s/${FASTLY_TEST_SERVICE_ID}/${DEFAULT_TEST_SERVICE_ID}/g" "$file"
done
