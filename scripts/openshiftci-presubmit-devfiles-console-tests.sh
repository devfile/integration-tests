#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/console $GOPATH/src/github.com/openshift/console
cp scripts/openshiftci-presubmit-devfiles-console-tests.sh $GOPATH/src/github.com/openshift/console/
cd $GOPATH/src/github.com/openshift/console
