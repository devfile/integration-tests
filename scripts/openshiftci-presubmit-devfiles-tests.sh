#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

git clone https://github.com/openshift/odo $GOPATH/src/github.com/openshift/odo
cp scripts/openshiftci-presubmit-devfiles-tests.sh $GOPATH/src/github.com/openshift/odo/scripts/
cd $GOPATH/src/github.com/openshift/odo

export CI="openshift"
make configure-installer-tests-cluster
make bin
mkdir -p $GOPATH/bin
make goget-ginkgo
export PATH="$PATH:$(pwd):$GOPATH/bin"
export ARTIFACTS_DIR="/tmp/artifacts"
export CUSTOM_HOMEDIR=$ARTIFACTS_DIR

# Copy kubeconfig to temporary kubeconfig file
# Read and Write permission to temporary kubeconfig file
TMP_DIR=$(mktemp -d)
cp $KUBECONFIG $TMP_DIR/kubeconfig
chmod 640 $TMP_DIR/kubeconfig
export KUBECONFIG=$TMP_DIR/kubeconfig

# Login as developer
odo login -u developer -p password@123

# Check login user name for debugging purpose
oc whoami

make test-integration-devfile

cp -r reports $ARTIFACTS_DIR

odo logout
