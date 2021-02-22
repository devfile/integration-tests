#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/console $GOPATH/src/github.com/openshift/console

cp scripts/console/test-cypress.sh $GOPATH/src/github.com/openshift/console/
cp scripts/console/test-prow-e2e.sh $GOPATH/src/github.com/openshift/console/

cp scripts/console/frontend/package.json $GOPATH/src/github.com/openshift/console/frontend/
cp scripts/console/frontend/packages/dev-console/integration-tests/features/addFlow/create-from-devfile.feature $GOPATH/src/github.com/openshift/console/frontend/packages/dev-console/integration-tests/features/addFlow/

cd $GOPATH/src/github.com/openshift/console
