#!/bin/bash
set -e

FIXTURESDIR="$(pwd)/fastly/fixtures/"
DEFAULT_SERVICE_ID="7i6HN3TK9wS159v2gPAZ8A"

if [[ -z $FASTLY_TEST_SERVICE_ID && -z $1 ]]; then
  echo "You must supply a service ID as either an argument or by setting \$FASTLY_TEST_SERVICE_ID"
  exit
fi

if [[ -z ${FASTLY_TEST_SERVICE_ID} ]]; then
  FASTLY_TEST_SERVICE_ID=$1
fi

echo "Searching fixtures for ${FASTLY_TEST_SERVICE_ID}"
grep --recursive --files-with-matches "${FASTLY_TEST_SERVICE_ID}" "${FIXTURESDIR}" | xargs sed -i "s/${FASTLY_TEST_SERVICE_ID}/${DEFAULT_SERVICE_ID}/g"
