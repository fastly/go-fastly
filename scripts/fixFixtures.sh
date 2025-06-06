#!/bin/bash
set -e

FASTLY_TEST_RESOURCE_ID=${1}
DEFAULT_TEST_RESOURCE_ID=${2}

if [[ -z "${1}" ]]; then
  echo "You must supply a resource ID as the first argument"
  exit
fi

if [[ -z "${2}" ]]; then
  echo "You must supply a resource ID as the second argument"
  exit
fi

find . -type d -name "fixtures" -print0 | while IFS= read -r -d $'\0' current_fixtures_dir; do
  for file in $(grep --recursive --files-with-matches "${FASTLY_TEST_RESOURCE_ID}" "${current_fixtures_dir}")
  do
    sed -i.bak "s/${FASTLY_TEST_RESOURCE_ID}/${DEFAULT_TEST_RESOURCE_ID}/g" "$file" && rm "${file}.bak"
  done
done
