#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

export ODO_DIRPATH=$GOPATH/src/github.com/redhat-developer/odo

# refresh odo source files
rm -rf $ODO_DIRPATH

## shallow clone only the target branch or tag to test
## e.g. git clone --depth=1 --branch=v2.2.2 https://github.com/redhat-developer/odo $ODO_DIRPATH
git clone --depth=1 https://github.com/redhat-developer/odo $ODO_DIRPATH

# overwrite with devfile/integration-tests Makefile
cp ./Makefile $ODO_DIRPATH/Makefile

# overwrite with devfile/integration-tests/*
rm -rf $ODO_DIRPATH/tests/integration/*
cp -r tests/integration/devfile/* $ODO_DIRPATH/tests/integration
cd $ODO_DIRPATH

make bin

mkdir -p $GOPATH/bin

# add $GOPATH/bin if it's not set
# export PATH="$PATH:$GOPATH/bin"

export REPORTS_DIR=$ODO_DIRPATH/tests/reports

# clean test reports directory
rm -rf $REPORTS_DIR
mkdir $REPORTS_DIR

# Copy kubeconfig to temporary kubeconfig file
# Read and Write permission to temporary kubeconfig file
#TMP_DIR=$(mktemp -d)
#cp $KUBECONFIG $TMP_DIR/kubeconfig
#chmod 640 $TMP_DIR/kubeconfig
#export KUBECONFIG=$TMP_DIR/kubeconfig

# Login as developer
oc login -u developer -p developer

# Check login user name for debugging purpose
oc whoami

### Test options. Uncomment one of tests below

# 1. run all devfile integration tests.
make test-integration || error=true

# 2. run individual integration test for ODO commands
 make test-cmd-devfile-list
# make test-cmd-devfile-init
# make test-cmd-devfile-push
# make test-cmd-devfile-exec
# make test-cmd-devfile-status
# make test-cmd-devfile-watch
# make test-cmd-devfile-app
# make test-cmd-delete
# make test-cmd-devfile-registry
# make test-cmd-devfile-test
# make test-cmd-devfile-debug
# make test-cmd-devfile-storage
# make test-cmd-devfile-log
# make test-cmd-devfile-env
# make test-cmd-devfile-config
# make test-cmd-watch
# make test-cmd-debug
# make test-interactive

# 3. run end-to-end devfile test
# make test-e2e-devfile

### Optional:
# merge multiple number(depending on TEST_EXEC_NODES in Makefile) of junit_*.xml into a single file.
# jrm $REPORTS_DIR/junit_combined.xml "$REPORTS_DIR/junit*.xml"
# convert test results from junit*.xml into HTML format
# xunit-viewer -r $REPORTS_DIR/junit_combined.xml -o $REPORTS_DIR/junit_combined.html -b https://raw.githubusercontent.com/josephca/devfile-icon/main/docs/icons/2021_Devfile_logo_DevLoop_Icon.png -t "Devfile performance test"

if [ $error ]; then
    exit -1
fi

oc logout
