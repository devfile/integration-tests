#!/usr/bin/env bash
set -euo pipefail
cd frontend
yarn install

function generateReport {
  yarn run cypress-postreport
  if test -f ./packages/integration-tests-cypress/cypress-a11y-report.json; then
    yarn cypress-a11y-report
  fi
}
trap generateReport EXIT

while getopts p:s:h:l: flag
do
  case "${flag}" in
    p) pkg=${OPTARG};;
    s) spec=${OPTARG};;
    h) headless=${OPTARG};;
  esac
done

if [ -n "${headless-}" ] && [ -z "${pkg-}" ]; then
  yarn run test-cypress-devconsole-headless-create-from-devfile
  exit;
fi

