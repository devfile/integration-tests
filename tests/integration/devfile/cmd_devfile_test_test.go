package devfile

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/odo/tests/helper"
	"path/filepath"
)

var _ = Describe("odo devfile test command tests", func() {
	var cmpName string
	var sourcePath = "/projects"
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		cmpName = helper.RandString(6)
		helper.Chdir(commonVar.Context)
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	Context("Should show proper errors", func() {

		// used ";" as consolidating symbol as this spec covers multiple scenerios
		It("should show error if component is not pushed; should error out if a non-existent command or a command from wrong group is specified", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--project", commonVar.Project, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))

			output := helper.CmdShouldFail("odo", "test", "--context", commonVar.Context)
			Expect(output).To(ContainSubstring("error occurred while getting the pod: pod not found for the selector"))

			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)
			output = helper.CmdShouldFail("odo", "test", "--test-command", "invalidcmd", "--context", commonVar.Context)
			Expect(output).To(ContainSubstring("not found in the devfile"))

			output = helper.CmdShouldFail("odo", "test", "--test-command", "devrun", "--context", commonVar.Context)
			Expect(output).To(ContainSubstring("command devrun is of group run in devfile.yaml"))
		})

		It("should show error if no test group is defined", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)
			output := helper.CmdShouldFail("odo", "test", "--context", commonVar.Context)

			Expect(output).To(ContainSubstring("the command group of kind \"test\" is not found in the devfile"))
		})

		It("should show error if devfile has no default test command", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.ReplaceString("devfile.yaml", "isDefault: true", "")
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)
			output := helper.CmdShouldFail("odo", "test", "--context", commonVar.Context)
			Expect(output).To(ContainSubstring("there should be exactly one default command for command group test, currently there is no default command"))
		})

		It("should show error if devfile has multiple default test command", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-multiple-defaults.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--build-command", "firstbuild", "--run-command", "secondrun", "--context", commonVar.Context)
			output := helper.CmdShouldFail("odo", "test", "--context", commonVar.Context)
			Expect(output).To(ContainSubstring("there should be exactly one default command for command group test, currently there is more than one default command"))
		})

		It("should error out on devfile flag", func() {
			helper.CmdShouldFail("odo", "test", "--devfile", "invalid.yaml")
		})
	})

	Context("Should run test command successfully", func() {

		It("Should run test command successfully with only one default specified", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)

			output := helper.CmdShouldPass("odo", "test", "--context", commonVar.Context)
			helper.MatchAllInOutput(output, []string{"Executing test1 command", "mkdir test1"})

			podName := commonVar.CliRunner.GetRunningPodNameByComponent(cmpName, commonVar.Project)
			output = commonVar.CliRunner.ExecListDir(podName, commonVar.Project, sourcePath)
			Expect(output).To(ContainSubstring("test1"))
		})

		It("Should run test command successfully with test-command specified", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)

			output := helper.CmdShouldPass("odo", "test", "--test-command", "test2", "--context", commonVar.Context)
			helper.MatchAllInOutput(output, []string{"Executing test2 command", "mkdir test2"})

			podName := commonVar.CliRunner.GetRunningPodNameByComponent(cmpName, commonVar.Project)
			output = commonVar.CliRunner.ExecListDir(podName, commonVar.Project, sourcePath)
			Expect(output).To(ContainSubstring("test2"))
		})

		It("should run test command successfully with test-command specified if devfile has no default test command", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.ReplaceString("devfile.yaml", "isDefault: true", "")
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)
			output := helper.CmdShouldPass("odo", "test", "--test-command", "test2", "--context", commonVar.Context)
			helper.MatchAllInOutput(output, []string{"Executing test2 command", "mkdir test2"})

			podName := commonVar.CliRunner.GetRunningPodNameByComponent(cmpName, commonVar.Project)
			output = commonVar.CliRunner.ExecListDir(podName, commonVar.Project, sourcePath)
			Expect(output).To(ContainSubstring("test2"))
		})

		It("should run test command successfully with test-command specified if devfile has multiple default test command", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-multiple-defaults.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--build-command", "firstbuild", "--run-command", "secondrun", "--context", commonVar.Context)
			output := helper.CmdShouldPass("odo", "test", "--test-command", "test2", "--context", commonVar.Context)
			helper.MatchAllInOutput(output, []string{"Executing test2 command", "mkdir test2"})

			podName := commonVar.CliRunner.GetRunningPodNameByComponent(cmpName, commonVar.Project)
			output = commonVar.CliRunner.ExecListDir(podName, commonVar.Project, sourcePath)
			Expect(output).To(ContainSubstring("test2"))
		})

		It("Should run composite test command successfully", func() {
			helper.CmdShouldPass("odo", "create", "nodejs", "--context", commonVar.Context, cmpName)

			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-with-testgroup.yaml"), filepath.Join(commonVar.Context, "devfile.yaml"))
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)

			output := helper.CmdShouldPass("odo", "test", "--test-command", "compositetest", "--context", commonVar.Context)
			helper.MatchAllInOutput(output, []string{"Executing test1 command", "mkdir test1", "Executing test2 command", "mkdir test2"})

			podName := commonVar.CliRunner.GetRunningPodNameByComponent(cmpName, commonVar.Project)
			output = commonVar.CliRunner.ExecListDir(podName, commonVar.Project, sourcePath)
			helper.MatchAllInOutput(output, []string{"test1", "test2"})
		})
	})

})
