package devfile

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/odo/tests/helper"
)

var _ = Describe("odo devfile registry command tests", func() {
	const registryName string = "TestDevfileRegistry"
	// Use staging OCI-based registry for tests to avoid a potential overload
	const addRegistryURL string = "https://registry.stage.devfile.io"

	const updateRegistryURL string = "http://www.example.com/update"
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		helper.Chdir(commonVar.Context)
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	Context("When executing registry list", func() {
		PIt("Should list all default registries", func() {
			output := helper.CmdShouldPass("odo", "registry", "list")
			helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
		})

		Measure("Should list all default registries", func(b Benchmarker) {
			runtime := b.Time("========== Command: odo registry list ==========", func() {
				output := helper.CmdShouldPass("odo", "registry", "list")
				helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 200), "odo registry list command should take less than 200 ms.")
			b.RecordValueWithPrecision("========== Execution time in ms ==========", float64(runtime.Milliseconds()), "ms", 2)
		}, 10)

		PIt("Should list all default registries with json", func() {
			output := helper.CmdShouldPass("odo", "registry", "list", "-o", "json")
			helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
		})

		Measure("Should list all default registries with json", func(b Benchmarker) {
			runtime := b.Time("========== Command: odo registry list -o json ==========", func() {
				output := helper.CmdShouldPass("odo", "registry", "list", "-o", "json")
				helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 200), "odo registry list -o json command should take less than 200 ms.")
			b.RecordValueWithPrecision("========== Execution time in ms ==========", float64(runtime.Milliseconds()), "ms", 2)
		}, 10)

		It("Should fail with an error with no registries", func() {
			helper.CmdShouldPass("odo", "registry", "delete", "DefaultDevfileRegistry", "-f")
			output := helper.CmdShouldFail("odo", "registry", "list")
			helper.MatchAllInOutput(output, []string{"No devfile registries added to the configuration. Refer `odo registry add -h` to add one"})
		})
	})

	Context("When executing registry commands with the registry is not present", func() {
		It("Should successfully add the registry", func() {
			helper.CmdShouldPass("odo", "registry", "add", registryName, addRegistryURL)
			output := helper.CmdShouldPass("odo", "registry", "list")
			helper.MatchAllInOutput(output, []string{registryName, addRegistryURL})
			helper.CmdShouldPass("odo", "create", "nodejs", "--registry", registryName)
			helper.CmdShouldPass("odo", "registry", "delete", registryName, "-f")
		})

		It("Should fail to update the registry", func() {
			helper.CmdShouldFail("odo", "registry", "update", registryName, updateRegistryURL, "-f")
		})

		It("Should fail to delete the registry", func() {
			helper.CmdShouldFail("odo", "registry", "delete", registryName, "-f")
		})
	})

	Context("When executing registry commands with the registry is present", func() {
		It("Should fail to add the registry", func() {
			helper.CmdShouldPass("odo", "registry", "add", registryName, addRegistryURL)
			helper.CmdShouldFail("odo", "registry", "add", registryName, addRegistryURL)
			helper.CmdShouldPass("odo", "registry", "delete", registryName, "-f")
		})

		It("Should successfully update the registry", func() {
			helper.CmdShouldPass("odo", "registry", "add", registryName, addRegistryURL)
			helper.CmdShouldPass("odo", "registry", "update", registryName, updateRegistryURL, "-f")
			output := helper.CmdShouldPass("odo", "registry", "list")
			helper.MatchAllInOutput(output, []string{registryName, updateRegistryURL})
			helper.CmdShouldPass("odo", "registry", "delete", registryName, "-f")
		})

		It("Should successfully delete the registry", func() {
			helper.CmdShouldPass("odo", "registry", "add", registryName, addRegistryURL)
			helper.CmdShouldPass("odo", "registry", "delete", registryName, "-f")
			helper.CmdShouldFail("odo", "create", "java-maven", "--registry", registryName)
		})

	})
})
