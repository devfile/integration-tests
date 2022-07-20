#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

export ODO_DIRPATH=$GOPATH/src/github.com/redhat-developer/odo

git clone --depth=1 https://github.com/redhat-developer/odo $ODO_DIRPATH
cp scripts/openshiftci-presubmit-devfiles-odo-all-tests.sh $ODO_DIRPATH/scripts/

mkdir $ODO_DIRPATH/tests/devfile-tests
cp $ODO_DIRPATH/tests/integration/cmd_devfile*.go $ODO_DIRPATH/tests/devfile-tests
rm -rf $ODO_DIRPATH/tests/integration/*
cp $ODO_DIRPATH/tests/devfile-tests/cmd_devfile*.go $ODO_DIRPATH/tests/integration
rm -rf $ODO_DIRPATH/tests/devfile-tests
cd $ODO_DIRPATH

# Run performance tests on top of integration tests
# sed -i 's/-randomizeAllSpecs/--noisyPendings=false/g' $ODO_DIRPATH/Makefile

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

# Integration tests
make test-integration || error=true

if [ $error ]; then
    exit -1
fi

oc logout
