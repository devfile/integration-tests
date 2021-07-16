#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone --depth=1 https://github.com/openshift/console $GOPATH/src/github.com/openshift/console

cp scripts/console/test-cypress.sh $GOPATH/src/github.com/openshift/console/
cp scripts/console/test-prow-e2e.sh $GOPATH/src/github.com/openshift/console/

sed -i'' '/test-cypress-dev-console-headless/a \\t"test-cypress-dev-console-headless-create-from-devfile": "cd packages/dev-console/integration-tests && node --max-old-space-size=4096 ../../../node_modules/.bin/cypress run --env openshift=true --browser ${BRIDGE_E2E_BROWSER_NAME:=chrome} --headless --spec \\"features/addFlow/create-from-devfile.feature\\";"\,' $GOPATH/src/github.com/openshift/console/frontend/package.json

cp scripts/console/frontend/packages/dev-console/integration-tests/features/addFlow/create-from-devfile.feature $GOPATH/src/github.com/openshift/console/frontend/packages/dev-console/integration-tests/features/addFlow/

cd $GOPATH/src/github.com/openshift/console
