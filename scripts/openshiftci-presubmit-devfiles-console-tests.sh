#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/console $GOPATH/src/github.com/openshift/console
cp scripts/console/test-cypress.sh $GOPATH/src/github.com/openshift/console/
cp scripts/console/test-prow-e2e.sh $GOPATH/src/github.com/openshift/console/
cp scripts/console/frontend/package.json $GOPATH/src/github.com/openshift/console/frontend/
cd $GOPATH/src/github.com/openshift/console
