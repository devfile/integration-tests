#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

ARCH=$(uname -m)
export CI="openshift"
if [ "${ARCH}" == "s390x" ]; then
    make configure-installer-tests-cluster-s390x
elif [ "${ARCH}" == "ppc64le" ]; then
    make configure-installer-tests-cluster-ppc64le
else
    make configure-installer-tests-cluster
fi
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
odo login -u developer -p developer

# Check login user name for debugging purpose
oc whoami

if [ "${ARCH}" == "s390x" ]; then
    echo "No integration tests for ${ARCH}"
elif  [ "${ARCH}" == "ppc64le" ]; then
    echo "No integration tests for ${ARCH}"
else
    # Integration tests
    echo "Run devfile integration tests"
    make test-integration-devfile
fi

odo logout
