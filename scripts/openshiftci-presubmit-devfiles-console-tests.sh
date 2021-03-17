#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/console $GOPATH/src/github.com/openshift/console

cp scripts/console/test-cypress.sh $GOPATH/src/github.com/openshift/console/
cp scripts/console/test-prow-e2e.sh $GOPATH/src/github.com/openshift/console/

sed -i.bak 's|"test-cypress-devconsole": "cd packages/dev-console/integration-tests && ../../../node_modules/.bin/cypress open --env openshift=true",|"test-cypress-devconsole-headless-create-from-devfile": "cd packages/dev-console/integration-tests && node --max-old-space-size=4096 ../../../node_modules/.bin/cypress run --env openshift=true --browser ${BRIDGE_E2E_BROWSER_NAME:=chrome} --headless --spec \"features/addFlow/create-from-devfile.feature\";",|' $GOPATH/src/github.com/openshift/console/frontend/package.json
rm $GOPATH/src/github.com/openshift/console/frontend/package.json.bak

cp scripts/console/frontend/packages/dev-console/integration-tests/features/addFlow/create-from-devfile.feature $GOPATH/src/github.com/openshift/console/frontend/packages/dev-console/integration-tests/features/addFlow/

cd $GOPATH/src/github.com/openshift/console
