#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

export ODO_DIRPATH=$GOPATH/src/github.com/redhat-developer/odo

git clone --depth=1 https://github.com/redhat-developer/odo $ODO_DIRPATH
cp scripts/openshiftci-presubmit-devfiles-odo-tests.sh $ODO_DIRPATH/scripts/

mkdir $ODO_DIRPATH/tests/devfile-tests
cp $ODO_DIRPATH/tests/integration/cmd_*.go $ODO_DIRPATH/tests/devfile-tests
rm -rf $ODO_DIRPATH/tests/integration/*
cp $ODO_DIRPATH/tests/devfile-tests/cmd_*.go $ODO_DIRPATH/tests/integration
rm -rf $ODO_DIRPATH/tests/devfile-tests
cd $ODO_DIRPATH

# Update with the latest devfile library for tests
go get -d github.com/devfile/library@main
go mod tidy
go mod vendor

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

#cp -r $ODO_DIRPATH/tests/integration/reports $ARTIFACT_DIR

oc logout
