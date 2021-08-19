#!/usr/bin/env bash
set -euo pipefail
cd frontend
yarn config set ignore-engines true
yarn install

function generateReport {
  yarn run cypress-postreport
  if test -f ./packages/integration-tests-cypress/cypress-a11y-report.json; then
    yarn cypress-a11y-report
  fi
}
trap generateReport EXIT

yarn run test-cypress-dev-console-headless-create-from-devfile

