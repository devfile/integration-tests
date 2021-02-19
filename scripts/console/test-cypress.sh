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

yarn_script="cd packages/dev-console/integration-tests && node --max-old-space-size=4096 ../../../node_modules/.bin/cypress run --env openshift=true --browser ${BRIDGE_E2E_BROWSER_NAME:=chrome} --headless --spec \"features/addFlow/create-from-devfile.feature\";"

yarn run $yarn_script
