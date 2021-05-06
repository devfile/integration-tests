package devfile

import (
	"testing"

	"github.com/openshift/odo/tests/helper"
)

func TestDevfiles(t *testing.T) {
	helper.RunTestSpecs(t, "Devfile Suite")
}

var _ = BeforeSuite(func() {
	const registryName string = "DefaultDevfileRegistry"
	// Use staging OCI-based registry for tests to avoid a potential overload
	const addRegistryURL string = "https://registry.stage.devfile.io"

	helper.CmdShouldPass("odo", "registry", "delete", "DefaultDevfileRegistry", "-f")
	helper.CmdShouldPass("odo", "registry", "add", registryName, addRegistryURL)
})
