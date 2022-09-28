//
// Copyright 2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package devfile

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"github.com/openshift/odo/tests/helper"
)

var _ = Describe("odo devfile registry command tests", func() {
	const registryName string = "RegistryName"
	// Use staging OCI-based registry for tests to avoid overload
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
		It("Should list all default registries", func() {
			output := helper.Cmd("odo", "registry", "list").ShouldPass().Out()
			helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
		})

		Measure("Should list all default registries", func(b Benchmarker) {
			runtime := b.Time("========== Command: odo registry list ==========", func() {
				output := helper.Cmd("odo", "registry", "list").ShouldPass().Out()
				helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 200), "odo registry list command should take less than 200 ms.")
			b.RecordValueWithPrecision("========== Execution time in ms ==========", float64(runtime.Milliseconds()), "ms", 2)
		}, 10)

		It("Should list all default registries with json", func() {
			output := helper.Cmd("odo", "registry", "list", "-o", "json").ShouldPass().Out()
			helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
		})

		Measure("Should list all default registries with json", func(b Benchmarker) {
			runtime := b.Time("========== Command: odo registry list -o json ==========", func() {
				output := helper.Cmd("odo", "registry", "list", "-o", "json").ShouldPass().Out()
				helper.MatchAllInOutput(output, []string{"DefaultDevfileRegistry"})
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 200), "odo registry list -o json command should take less than 200 ms.")
			b.RecordValueWithPrecision("========== Execution time in ms ==========", float64(runtime.Milliseconds()), "ms", 2)
		}, 10)

		It("Should fail with an error with no registries", func() {
			helper.Cmd("odo", "registry", "delete", "DefaultDevfileRegistry", "-f").ShouldPass()
			output := helper.Cmd("odo", "registry", "list").ShouldFail().Err()
			helper.MatchAllInOutput(output, []string{"No devfile registries added to the configuration. Refer `odo registry add -h` to add one"})

		})

	})

	Context("When executing registry commands with the registry is not present", func() {
		It("Should successfully add the registry", func() {
			helper.Cmd("odo", "registry", "add", registryName, addRegistryURL).ShouldPass()
			output := helper.Cmd("odo", "registry", "list").ShouldPass().Out()
			helper.MatchAllInOutput(output, []string{registryName, addRegistryURL})
			helper.Cmd("odo", "create", "nodejs", "--registry", registryName).ShouldPass()
			helper.Cmd("odo", "registry", "delete", registryName, "-f").ShouldPass()
		})

		It("Should fail to update the registry", func() {
			helper.Cmd("odo", "registry", "update", registryName, updateRegistryURL, "-f").ShouldFail()
		})

		It("Should fail to delete the registry", func() {
			helper.Cmd("odo", "registry", "delete", registryName, "-f").ShouldFail()
		})
	})

	Context("When executing registry commands with the registry is present", func() {
		It("Should fail to add the registry", func() {
			helper.Cmd("odo", "registry", "add", registryName, addRegistryURL).ShouldPass()
			helper.Cmd("odo", "registry", "add", registryName, addRegistryURL).ShouldFail()
			helper.Cmd("odo", "registry", "delete", registryName, "-f").ShouldPass()
		})

		It("Should successfully update the registry", func() {
			helper.Cmd("odo", "registry", "add", registryName, addRegistryURL).ShouldPass()
			helper.Cmd("odo", "registry", "update", registryName, updateRegistryURL, "-f").ShouldPass()
			output := helper.Cmd("odo", "registry", "list").ShouldPass().Out()
			helper.MatchAllInOutput(output, []string{registryName, updateRegistryURL})
			helper.Cmd("odo", "registry", "delete", registryName, "-f").ShouldPass()
		})

		It("Should successfully delete the registry", func() {
			helper.Cmd("odo", "registry", "add", registryName, addRegistryURL).ShouldPass()
			helper.Cmd("odo", "registry", "delete", registryName, "-f").ShouldPass()
			helper.Cmd("odo", "create", "java-maven", "--registry", registryName).ShouldFail()
		})
	})

	Context("when working with git based registries", func() {
		var deprecated, docLink string
		JustBeforeEach(func() {
			deprecated = "Deprecated"
			docLink = "https://github.com/openshift/odo/tree/main/docs/public/git-registry-deprecation.adoc"
		})
		It("should show deprecation warning when the git based registry is used", func() {

			outstr, errstr := helper.Cmd("odo", "registry", "add", "RegistryFromGitHub", "https://github.com/odo-devfiles/registry").ShouldPass().OutAndErr()
			co := fmt.Sprintln(outstr, errstr)
			helper.MatchAllInOutput(co, []string{deprecated, docLink})
			outstr, errstr = helper.Cmd("odo", "registry", "list").ShouldPass().OutAndErr()
			co = fmt.Sprintln(outstr, errstr)
			helper.MatchAllInOutput(co, []string{deprecated, docLink})
			outstr, errstr = helper.Cmd("odo", "create", "nodejs", "--registry", "RegistryFromGitHub").ShouldPass().OutAndErr()
			co = fmt.Sprintln(outstr, errstr)
			helper.MatchAllInOutput(co, []string{deprecated, docLink})
		})
		It("should not show deprecation warning if non-git-based registry is used", func() {
			out, err := helper.Cmd("odo", "registry", "list").ShouldPass().OutAndErr()
			helper.DontMatchAllInOutput(fmt.Sprintln(out, err), []string{deprecated, docLink})
			helper.Cmd("odo", "registry", "add", "RegistryFromGitHub", "https://github.com/odo-devfiles/registry").ShouldPass()
			out, err = helper.Cmd("odo", "create", "nodejs", "--registry", "DefaultDevfileRegistry").ShouldPass().OutAndErr()
			helper.DontMatchAllInOutput(fmt.Sprintln(out, err), []string{deprecated, docLink})
		})
	})
})
