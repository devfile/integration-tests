#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/odo $GOPATH/src/github.com/openshift/odo
cp scripts/openshiftci-presubmit-devfiles-odo-tests.sh $GOPATH/src/github.com/openshift/odo/scripts/
cd $GOPATH/src/github.com/openshift/odo

export CI="openshift"
make configure-installer-tests-cluster
make bin
mkdir -p $GOPATH/bin
make goget-ginkgo
export PATH="$PATH:$(pwd):$GOPATH/bin"
export CUSTOM_HOMEDIR=$ARTIFACT_DIR

# Copy kubeconfig to temporary kubeconfig file
# Read and Write permission to temporary kubeconfig file
TMP_DIR=$(mktemp -d)
cp $KUBECONFIG $TMP_DIR/kubeconfig
chmod 640 $TMP_DIR/kubeconfig
export KUBECONFIG=$TMP_DIR/kubeconfig

# Login as developer
oc login -u developer -p password@123

# Check login user name for debugging purpose
oc whoami

make test-integration-devfile

cp -r $GOPATH/src/github.com/openshift/odo/tests/integration/reports $ARTIFACT_DIR

oc logout
