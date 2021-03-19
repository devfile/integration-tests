#!/usr/bin/env bash
set -euo pipefail
cd frontend
# yarn config set ignore-engines true
rm -rf node_modules/ package-lock.json yarn.lock
yarn install

function generateReport {
  yarn run cypress-postreport
  if test -f ./packages/integration-tests-cypress/cypress-a11y-report.json; then
    yarn cypress-a11y-report
  fi
}
trap generateReport EXIT

yarn run test-cypress-devconsole-headless-create-from-devfile

