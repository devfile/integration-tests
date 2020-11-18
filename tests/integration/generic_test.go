package integration

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/odo/tests/helper"
)

var _ = Describe("odo generic", func() {
	var testPHPGitURL = "https://github.com/appuio/example-php-sti-helloworld"
	var testNodejsGitURL = "https://github.com/sclorg/nodejs-ex"
	var testLongURLName = "long-url-name-long-url-name-long-url-name-long-url-name-long-url-name"

	// TODO: A neater way to provide odo path. Currently we assume \
	// odo and oc in $PATH already
	var oc helper.OcRunner
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		oc = helper.NewOcRunner("oc")
		commonVar = helper.CommonBeforeEach()
	})

	// Clean up after the test
	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	Context("Check the help usage for odo", func() {

		It("Makes sure that we have the long-description when running odo and we dont error", func() {
			output := helper.CmdShouldPass("odo")
			Expect(output).To(ContainSubstring("To see a full list of commands, run 'odo --help'"))
		})

		It("Make sure we have the full description when performing odo --help", func() {
			output := helper.CmdShouldPass("odo", "--help")
			Expect(output).To(ContainSubstring("Use \"odo [command] --help\" for more information about a command."))
		})

		It("Fail when entering an incorrect name for a component", func() {
			output := helper.CmdShouldFail("odo", "component", "foobar")
			Expect(output).To(ContainSubstring("Subcommand not found, use one of the available commands"))
		})

		It("Fail with showing help only once for incorrect command", func() {
			output := helper.CmdShouldFail("odo", "hello")
			Expect(strings.Count(output, "odo [flags]")).Should(Equal(1))
		})

	})

	Context("When executing catalog list without component directory", func() {
		It("should list all component catalogs", func() {
			stdOut := helper.CmdShouldPass("odo", "catalog", "list", "components")
			helper.MatchAllInOutput(stdOut, []string{"dotnet", "nginx", "php", "ruby", "wildfly"})
		})

	})

	Context("check catalog component search functionality", func() {
		It("check that a component does not exist", func() {
			componentRandomName := helper.RandString(7)
			output := helper.CmdShouldFail("odo", "catalog", "search", "component", componentRandomName)
			Expect(output).To(ContainSubstring("no component matched the query: " + componentRandomName))
		})
	})

	// Test machine readable output
	Context("when creating project -o json", func() {
		var projectName string
		JustBeforeEach(func() {
			projectName = helper.RandString(6)
		})
		JustAfterEach(func() {
			helper.DeleteProject(projectName)
		})

		// odo project create foobar -o json
		It("should be able to create project and show output in json format", func() {
			actual := helper.CmdShouldPass("odo", "project", "create", projectName, "-o", "json")
			desired := fmt.Sprintf(`{"kind":"Project","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"%s","namespace":"%s","creationTimestamp":null},"message":"Project '%s' is ready for use"}`, projectName, projectName, projectName)
			Expect(desired).Should(MatchJSON(actual))
		})
	})

	Context("Creating same project twice with flag -o json", func() {
		var projectName string
		JustBeforeEach(func() {
			projectName = helper.RandString(6)
		})
		JustAfterEach(func() {
			helper.DeleteProject(projectName)
		})
		// odo project create foobar -o json (x2)
		It("should fail along with proper machine readable output", func() {
			helper.CmdShouldPass("odo", "project", "create", projectName)
			actual := helper.CmdShouldFail("odo", "project", "create", projectName, "-o", "json")
			desired := fmt.Sprintf(`{"kind":"Error","apiVersion":"odo.dev/v1alpha1","metadata":{"creationTimestamp":null},"message":"unable to create new project: unable to create new project %s: project.project.openshift.io \"%s\" already exists"}`, projectName, projectName)
			Expect(desired).Should(MatchJSON(actual))
		})
	})

	Context("Delete the project with flag -o json", func() {
		var projectName string
		JustBeforeEach(func() {
			projectName = helper.RandString(6)
		})

		// odo project delete foobar -o json
		It("should be able to delete project and show output in json format", func() {
			helper.CmdShouldPass("odo", "project", "create", projectName, "-o", "json")

			actual := helper.CmdShouldPass("odo", "project", "delete", projectName, "-o", "json")
			desired := fmt.Sprintf(`{"kind":"Project","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"%s","namespace":"%s","creationTimestamp":null},"message":"Deleted project : %s"}`, projectName, projectName, projectName)
			Expect(desired).Should(MatchJSON(actual))
		})
	})

	Context("creating component with an application and url", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})
		It("should create the component in default application", func() {
			helper.CmdShouldPass("odo", "create", "--s2i", "php", "testcmp", "--app", "e2e-xyzk", "--project", commonVar.Project, "--git", testPHPGitURL)
			helper.CmdShouldPass("odo", "config", "set", "Ports", "8080/TCP", "-f")
			helper.CmdShouldPass("odo", "push")
			oc.VerifyCmpName("testcmp", commonVar.Project)
			oc.VerifyAppNameOfComponent("testcmp", "e2e-xyzk", commonVar.Project)
			helper.CmdShouldPass("odo", "app", "delete", "e2e-xyzk", "-f")
		})
	})

	Context("Overwriting build timeout for git component", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})

		It("should pass to build component if the given build timeout is more than the default(300s) value", func() {
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project, "--git", testNodejsGitURL)
			helper.CmdShouldPass("odo", "preference", "set", "BuildTimeout", "600")
			buildTimeout := helper.GetPreferenceValue("BuildTimeout")
			helper.MatchAllInOutput(buildTimeout, []string{"600"})
			helper.CmdShouldPass("odo", "push")
		})

		It("should fail to build component if the given build timeout is pretty less(2s)", func() {
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project, "--git", testNodejsGitURL)
			helper.CmdShouldPass("odo", "preference", "set", "BuildTimeout", "2")
			buildTimeout := helper.GetPreferenceValue("BuildTimeout")
			helper.MatchAllInOutput(buildTimeout, []string{"2"})
			stdOut := helper.CmdShouldFail("odo", "push")
			helper.MatchAllInOutput(stdOut, []string{"Failed to create component", "timeout waiting for build"})
		})
	})

	Context("should list applications in other project", func() {
		It("should be able to create a php component with application created", func() {
			helper.CmdShouldPass("odo", "create", "--s2i", "php", "testcmp", "--app", "testing", "--project", commonVar.Project, "--ref", "master", "--git", testPHPGitURL, "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)
			currentProject := helper.CreateRandProject()
			currentAppNames := helper.CmdShouldPass("odo", "app", "list", "--project", currentProject)
			Expect(currentAppNames).To(ContainSubstring("There are no applications deployed in the project '" + currentProject + "'"))
			appNames := helper.CmdShouldPass("odo", "app", "list", "--project", commonVar.Project)
			Expect(appNames).To(ContainSubstring("testing"))
			helper.DeleteProject(currentProject)
		})
	})

	Context("when running odo push with flag --show-log", func() {
		It("should be able to push changes", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project, "--context", commonVar.Context)

			// Push the changes with --show-log
			getLogging := helper.CmdShouldPass("odo", "push", "--show-log", "--context", commonVar.Context)
			Expect(getLogging).To(ContainSubstring("Building component"))
		})
	})

	Context("deploying a component with a specific image name", func() {
		It("should deploy the component", func() {
			helper.CmdShouldPass("git", "clone", "https://github.com/openshift/nodejs-ex", commonVar.Context+"/nodejs-ex")
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs:latest", "testversioncmp", "--project", commonVar.Project, "--context", commonVar.Context+"/nodejs-ex")
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context+"/nodejs-ex")
			helper.CmdShouldPass("odo", "delete", "-f", "--context", commonVar.Context+"/nodejs-ex")
		})
	})

	Context("When deleting two project one after the other", func() {
		It("should be able to delete sequentially", func() {
			project1 := helper.CreateRandProject()
			project2 := helper.CreateRandProject()

			helper.DeleteProject(project2)
			helper.DeleteProject(project1)
		})
	})

	Context("When deleting three project one after the other in opposite order", func() {
		It("should be able to delete", func() {
			project1 := helper.CreateRandProject()
			project2 := helper.CreateRandProject()
			project3 := helper.CreateRandProject()

			helper.DeleteProject(project1)
			helper.DeleteProject(project2)
			helper.DeleteProject(project3)

		})
	})

	Context("when executing odo version command", func() {
		It("should show the version of odo major components including server login URL", func() {
			odoVersion := helper.CmdShouldPass("odo", "version")
			reOdoVersion := regexp.MustCompile(`^odo\s*v[0-9]+.[0-9]+.[0-9]+(?:-\w+)?\s*\(\w+\)`)
			odoVersionStringMatch := reOdoVersion.MatchString(odoVersion)
			rekubernetesVersion := regexp.MustCompile(`Kubernetes:\s*v[0-9]+.[0-9]+.[0-9]+((-\w+\.[0-9]+)?\+\w+)?`)
			kubernetesVersionStringMatch := rekubernetesVersion.MatchString(odoVersion)
			Expect(odoVersionStringMatch).Should(BeTrue())
			Expect(kubernetesVersionStringMatch).Should(BeTrue())
			serverURL := oc.GetCurrentServerURL()
			Expect(odoVersion).Should(ContainSubstring("Server: " + serverURL))
		})
	})

	Context("prevent the user from creating invalid URLs", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})
		It("should not allow creating a URL with long name", func() {
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", "--project", commonVar.Project)
			stdOut := helper.CmdShouldFail("odo", "url", "create", testLongURLName, "--port", "8080")
			Expect(stdOut).To(ContainSubstring("must be shorter than 63 characters"))
		})
	})

	Context("when component's deployment config is deleted with oc", func() {
		var componentRandomName string

		JustBeforeEach(func() {
			componentRandomName = helper.RandString(6)
			helper.Chdir(commonVar.Context)
		})

		It("should delete all OpenShift objects except the component's imagestream", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", componentRandomName, "--project", commonVar.Project)
			helper.CmdShouldPass("odo", "push")

			// Delete the deployment config using oc delete
			dc := oc.GetDcName(componentRandomName, commonVar.Project)
			helper.CmdShouldPass("oc", "delete", "--wait", "dc", dc, "--namespace", commonVar.Project)

			// insert sleep because it takes a few seconds to delete *all*
			// objects owned by DC but we should be able to check if a service
			// got deleted in a second.
			time.Sleep(1 * time.Second)

			// now check if the service owned by the DC exists. Service name is
			// same as DC name for a given component.
			stdOut := helper.CmdShouldFail("oc", "get", "svc", dc, "--namespace", commonVar.Project)
			Expect(stdOut).To(ContainSubstring("NotFound"))

			// ensure that the image stream still exists
			helper.CmdShouldPass("oc", "get", "is", dc, "--namespace", commonVar.Project)
		})
	})

	Context("When using cpu or memory flag with odo create", func() {
		cmpName := "nodejs"

		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})

		It("should not allow using any memory or cpu flag", func() {

			cases := []struct {
				paramName  string
				paramValue string
			}{
				{
					paramName:  "cpu",
					paramValue: "0.4",
				},
				{
					paramName:  "mincpu",
					paramValue: "0.2",
				},
				{
					paramName:  "maxcpu",
					paramValue: "0.4",
				},
				{
					paramName:  "memory",
					paramValue: "200Mi",
				},
				{
					paramName:  "minmemory",
					paramValue: "100Mi",
				},
				{
					paramName:  "maxmemory",
					paramValue: "200Mi",
				},
			}
			for _, testCase := range cases {
				helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
				output := helper.CmdShouldFail("odo", "component", "create", "--s2i", "nodejs", cmpName, "--project", commonVar.Project, "--context", commonVar.Context, "--"+testCase.paramName, testCase.paramValue)
				Expect(output).To(ContainSubstring("unknown flag: --" + testCase.paramName))
			}
		})
	})

})
