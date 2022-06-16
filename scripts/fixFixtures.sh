#!/bin/bash
set -e

FASTLY_TEST_SERVICE_ID=$1
DEFAULT_TEST_SERVICE_ID="7i6HN3TK9wS159v2gPAZ8A"
DEFAULT_TEST_USER_ID="4tKBSuFhNEiIpNDxmmVydt"
FIXTURESDIR="$(pwd)/fastly/fixtures/"

if [[ -z ${FASTLY_TEST_SERVICE_ID} && -z $1 ]]; then
  echo "You must supply a service ID as the first argument"
  exit
fi

if [[ -z ${FASTLY_TEST_USER_ID} && -z $2 ]]; then
  echo "You must supply a user ID as the second argument"
  exit
fi

for file in $(grep --recursive --files-with-matches "${FASTLY_TEST_SERVICE_ID}" "${FIXTURESDIR}")
do
  sed -i.bak "s/${FASTLY_TEST_SERVICE_ID}/${DEFAULT_TEST_SERVICE_ID}/g" "$file" && rm "${file}.bak"
done

for file in $(grep --recursive --files-with-matches "${FASTLY_TEST_USER_ID}" "${FIXTURESDIR}")
do
  sed -i.bak "s/${FASTLY_TEST_USER_ID}/${DEFAULT_TEST_USER_ID}/g" "$file" && rm "${file}.bak"
done
