package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/odo/tests/helper"
)

func componentTests(args ...string) {
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

	Context("Generic machine readable output tests", func() {

		It("Command should fail if json is non-existent for a command", func() {
			output := helper.CmdShouldFail("odo", "version", "-o", "json")
			Expect(output).To(ContainSubstring("Machine readable output is not yet implemented for this command"))
		})

		It("Help for odo version should not contain machine output", func() {
			output := helper.CmdShouldPass("odo", "version", "--help")
			Expect(output).NotTo(ContainSubstring("Specify output format, supported format: json"))
		})

	})

	Context("Creating component", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})

		It("should create but not list component even in new project with --project and --context at the same time", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--git", "https://github.com/openshift/nodejs-ex", "--project", commonVar.Project, "--context", commonVar.Context, "--app", "testing")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context, "-v4")...)
			projectList := helper.CmdShouldPass("odo", "project", "list")
			Expect(projectList).To(ContainSubstring(commonVar.Project))
			helper.CmdShouldFail("odo", "list", "--project", commonVar.Project, "--context", commonVar.Context)
		})

		It("Without an application should create one", func() {
			componentName := helper.RandString(6)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, componentName, "--ref", "master", "--git", "https://github.com/openshift/nodejs-ex")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+componentName, "Application,app")
			helper.CmdShouldPass("odo", append(args, "push")...)
			appName := helper.CmdShouldPass("odo", "app", "list")
			Expect(appName).ToNot(BeEmpty())

			// checking if application name is set to "app"
			applicationName := helper.GetConfigValue("Application")
			Expect(applicationName).To(Equal("app"))

			// clean up
			helper.CmdShouldPass("odo", "app", "delete", "app", "-f")
			helper.CmdShouldFail("odo", "app", "delete", "app", "-f")
			helper.CmdShouldFail("odo", append(args, "delete", componentName, "-f")...)

		})

		It("should create default named component when passed same context differently", func() {
			dir := filepath.Base(commonVar.Context)
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "--context", ".", "--app", "testing")...)
			componentName := helper.GetConfigValueWithContext("Name", commonVar.Context)
			Expect(componentName).To(ContainSubstring("nodejs-" + dir))
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+componentName, "Application,testing")
			helper.DeleteDir(filepath.Join(commonVar.Context, ".odo"))
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "--context", commonVar.Context, "--app", "testing")...)
			newComponentName := helper.GetConfigValueWithContext("Name", commonVar.Context)
			Expect(newComponentName).To(ContainSubstring("nodejs-" + dir))
		})

		It("should show an error when ref flag is provided with sources except git", func() {
			outputErr := helper.CmdShouldFail("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "cmp-git", "--ref", "test")...)
			Expect(outputErr).To(ContainSubstring("the --ref flag is only valid for --git flag"))
		})

		It("create component twice fails from same directory", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project)...)
			output := helper.CmdShouldFail("odo", append(args, "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project)...)
			Expect(output).To(ContainSubstring("this directory already contains a component"))
		})

		It("should list out component in json format along with path flag", func() {
			var contextPath string
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,nodejs", "Application,app")
			if runtime.GOOS == "windows" {
				contextPath = strings.Replace(strings.TrimSpace(commonVar.Context), "\\", "\\\\", -1)
			} else {
				contextPath = strings.TrimSpace(commonVar.Context)
			}
			// this orders the json
			desired, err := helper.Unindented(fmt.Sprintf(`{"kind":"Component","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"nodejs","namespace":"%s","creationTimestamp":null},"spec":{"app":"app","type":"nodejs","sourceType": "local","ports":["8080/TCP"]},"status":{"context":"%s","state":"Not Pushed"}}`, commonVar.Project, contextPath))
			Expect(err).Should(BeNil())

			actual, err := helper.Unindented(helper.CmdShouldPass("odo", append(args, "list", "-o", "json", "--path", filepath.Dir(commonVar.Context))...))
			Expect(err).Should(BeNil())
			// since the tests are run parallel, there might be many odo component directories in the root folder
			// so we only check for the presence of the current one
			Expect(actual).Should(ContainSubstring(desired))
		})

		It("should list out pushed components of different projects in json format along with path flag", func() {
			var contextPath string
			var contextPath2 string
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "nodejs", "--project", commonVar.Project)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,nodejs", "Application,app")
			helper.CmdShouldPass("odo", append(args, "push")...)

			project2 := helper.CreateRandProject()
			context2 := helper.CreateNewContext()
			helper.Chdir(context2)
			helper.CopyExample(filepath.Join("source", "python"), context2)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "python", "python", "--project", project2)...)
			helper.ValidateLocalCmpExist(context2, "Type,python", "Name,python", "Application,app")
			helper.CmdShouldPass("odo", append(args, "push")...)

			if runtime.GOOS == "windows" {
				contextPath = strings.Replace(strings.TrimSpace(commonVar.Context), "\\", "\\\\", -1)
				contextPath2 = strings.Replace(strings.TrimSpace(context2), "\\", "\\\\", -1)
			} else {
				contextPath = strings.TrimSpace(commonVar.Context)
				contextPath2 = strings.TrimSpace(context2)
			}

			actual, err := helper.Unindented(helper.CmdShouldPass("odo", append(args, "list", "-o", "json", "--path", filepath.Dir(commonVar.Context))...))
			Expect(err).Should(BeNil())
			helper.Chdir(commonVar.Context)
			helper.DeleteDir(context2)
			helper.DeleteProject(project2)
			// this orders the json
			expected, err := helper.Unindented(fmt.Sprintf(`{"kind":"Component","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"nodejs","namespace":"%s","creationTimestamp":null},"spec":{"app":"app","type":"nodejs","sourceType": "local","ports":["8080/TCP"]},"status":{"context":"%s","state":"Pushed"}}`, commonVar.Project, contextPath))
			Expect(err).Should(BeNil())
			Expect(actual).Should(ContainSubstring(expected))
			// this orders the json
			expected, err = helper.Unindented(fmt.Sprintf(`{"kind":"Component","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"python","namespace":"%s","creationTimestamp":null},"spec":{"app":"app","type":"python","sourceType": "local","ports":["8080/TCP"]},"status":{"context":"%s","state":"Pushed"}}`, project2, contextPath2))
			Expect(err).Should(BeNil())
			Expect(actual).Should(ContainSubstring(expected))

		})

		It("should list the component", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			cmpList := helper.CmdShouldPass("odo", append(args, "list", "--project", commonVar.Project)...)
			Expect(cmpList).To(ContainSubstring("cmp-git"))
			actualCompListJSON := helper.CmdShouldPass("odo", append(args, "list", "--project", commonVar.Project, "-o", "json")...)
			desiredCompListJSON := fmt.Sprintf(`{"kind":"List","apiVersion":"odo.dev/v1alpha1","metadata":{},"s2iComponents":[{"kind":"Component","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"cmp-git","namespace":"%s","creationTimestamp":null},"spec":{"app":"testing","type":"nodejs","source":"https://github.com/openshift/nodejs-ex","sourceType":"git","env":[{"name":"DEBUG_PORT","value":"5858"}]},"status":{"state":"Pushed"}}],"devfileComponents":[]}`, commonVar.Project)
			Expect(desiredCompListJSON).Should(MatchJSON(actualCompListJSON))
			cmpAllList := helper.CmdShouldPass("odo", append(args, "list", "--all-apps")...)
			Expect(cmpAllList).To(ContainSubstring("cmp-git"))
			helper.CmdShouldPass("odo", append(args, "delete", "cmp-git", "-f")...)
		})

		It("should list the component when it is not pushed", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing")
			cmpList := helper.CmdShouldPass("odo", append(args, "list", "--context", commonVar.Context)...)
			helper.MatchAllInOutput(cmpList, []string{"cmp-git", "Not Pushed"})
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
		})
		It("should list the state as unknown for disconnected cluster", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing")
			kubeconfigOrig := os.Getenv("KUBECONFIG")

			unset := func() {
				// KUBECONFIG defaults to ~/.kube/config so it can be empty in some cases.
				if kubeconfigOrig != "" {
					os.Setenv("KUBECONFIG", kubeconfigOrig)
				} else {
					os.Unsetenv("KUBECONFIG")
				}
			}

			os.Setenv("KUBECONFIG", "/no/such/path")

			defer unset()
			cmpList := helper.CmdShouldPass("odo", append(args, "list", "--context", commonVar.Context, "--v", "9")...)

			helper.MatchAllInOutput(cmpList, []string{"cmp-git", "Unknown"})
			unset()

			fmt.Printf("kubeconfig before delete %v", os.Getenv("KUBECONFIG"))
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
		})

		It("should describe the component when it is not pushed", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.CmdShouldPass("odo", "url", "create", "url-1", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "url", "create", "url-2", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "storage", "create", "storage-1", "--size", "1Gi", "--path", "/data1", "--context", commonVar.Context)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing", "URL,0,Name,url-1")
			cmpDescribe := helper.CmdShouldPass("odo", append(args, "describe", "--context", commonVar.Context)...)
			helper.MatchAllInOutput(cmpDescribe, []string{
				"cmp-git",
				"nodejs",
				"url-1",
				"url-2",
				"https://github.com/openshift/nodejs-ex",
				"storage-1",
			})

			cmpDescribeJSON, err := helper.Unindented(helper.CmdShouldPass("odo", append(args, "describe", "-o", "json", "--context", commonVar.Context)...))
			Expect(err).Should(BeNil())
			expected, err := helper.Unindented(`{"kind": "Component","apiVersion": "odo.dev/v1alpha1","metadata": {"name": "cmp-git","namespace": "` + commonVar.Project + `","creationTimestamp": null},"spec":{"app": "testing","type":"nodejs","source": "https://github.com/openshift/nodejs-ex","sourceType": "git","urls": {"kind": "List", "apiVersion": "odo.dev/v1alpha1", "metadata": {}, "items": [{"kind": "url", "apiVersion": "odo.dev/v1alpha1", "metadata": {"name": "url-1", "creationTimestamp": null}, "spec": {"port": 8080, "secure": false,  "kind": "route"}, "status": {"state": "Not Pushed"}}, {"kind": "url", "apiVersion": "odo.dev/v1alpha1", "metadata": {"name": "url-2", "creationTimestamp": null}, "spec": {"port": 8080, "secure": false,  "kind": "route"}, "status": {"state": "Not Pushed"}}]},"storages": {"kind": "List", "apiVersion": "odo.dev/v1alpha1", "metadata": {}, "items": [{"kind": "storage", "apiVersion": "odo.dev/v1alpha1", "metadata": {"name": "storage-1", "creationTimestamp": null}, "spec": {"size": "1Gi", "path": "/data1"}}]},"ports": ["8080/TCP"]},"status": {"state": "Not Pushed"}}`)
			Expect(err).Should(BeNil())
			Expect(cmpDescribeJSON).Should(MatchJSON(expected))

			// odo should describe not pushed component if component name is given.
			helper.CmdShouldPass("odo", append(args, "describe", "cmp-git", "--context", commonVar.Context)...)
			Expect(cmpDescribe).To(ContainSubstring("cmp-git"))

			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
		})

		It("checks that odo describe works for s2i component from a devfile directory", func() {
			newContext := path.Join(commonVar.Context, "newContext")
			helper.MakeDir(newContext)
			helper.Chdir(newContext)
			cmpName2 := helper.RandString(6)
			helper.CmdShouldPass("odo", "create", "--starter", "nodejs")
			context2 := helper.CreateNewContext()
			helper.CmdShouldPass("odo", "create", "--s2i", "nodejs", "--context", context2, cmpName2)
			output := helper.CmdShouldPass("odo", "describe", "--context", context2)
			Expect(output).To(ContainSubstring(fmt.Sprint("Component Name: ", cmpName2)))
			helper.Chdir(commonVar.OriginalWorkingDirectory)
			helper.DeleteDir(context2)
		})
		It("should describe not pushed component when it is created with json output", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			cmpDescribeJSON, err := helper.Unindented(helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--context", commonVar.Context, "--app", "testing", "-o", "json")...))
			Expect(err).Should(BeNil())
			expected, err := helper.Unindented(`{"kind": "Component","apiVersion": "odo.dev/v1alpha1","metadata": {"name": "cmp-git","namespace": "` + commonVar.Project + `","creationTimestamp": null},"spec":{"app": "testing","type":"nodejs","source": "file://./","sourceType": "local","ports": ["8080/TCP"]}, "status": {"state": "Not Pushed"}}`)
			Expect(err).Should(BeNil())
			Expect(cmpDescribeJSON).Should(MatchJSON(expected))
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
		})

		It("should describe pushed component when it is created with json output", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			cmpDescribeJSON, err := helper.Unindented(helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--context", commonVar.Context, "--app", "testing", "-o", "json", "--now")...))
			Expect(err).Should(BeNil())
			expected, err := helper.Unindented(`{"kind": "Component","apiVersion": "odo.dev/v1alpha1","metadata": {"name": "cmp-git","namespace": "` + commonVar.Project + `","creationTimestamp": null},"spec":{"app": "testing","type":"nodejs","sourceType": "local","env": [{"name": "DEBUG_PORT","value": "5858"}],"ports": ["8080/TCP"]}, "status": {"state": "Pushed"}}`)
			Expect(err).Should(BeNil())
			Expect(cmpDescribeJSON).Should(MatchJSON(expected))
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
		})

		It("should list the component in the same app when one is pushed and the other one is not pushed", func() {
			helper.Chdir(commonVar.OriginalWorkingDirectory)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,cmp-git", "Application,testing")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			context2 := helper.CreateNewContext()
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git-2", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", context2, "--app", "testing")...)
			helper.ValidateLocalCmpExist(context2, "Type,nodejs", "Name,cmp-git-2", "Application,testing")
			cmpList := helper.CmdShouldPass("odo", append(args, "list", "--context", context2)...)
			helper.MatchAllInOutput(cmpList, []string{"cmp-git", "cmp-git-2", "Not Pushed", "Pushed"})

			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--all", "--context", context2)...)
			helper.DeleteDir(context2)
		})

		It("should succeed listing catalog components", func() {
			// Since components catalog is constantly changing, we simply check to see if this command passes.. rather than checking the JSON each time.
			helper.CmdShouldPass("odo", "catalog", "list", "components", "-o", "json")
		})

		It("binary component should not fail when --context is not set", func() {
			oc.ImportJavaIS(commonVar.Project)
			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)
			// Was failing due to https://github.com/openshift/odo/issues/1969
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "java:8", "sb-jar-test", "--project",
				commonVar.Project, "--binary", filepath.Join(commonVar.Context, "sb.jar"))...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,java:8", "Name,sb-jar-test")
		})

		It("binary component should fail when --binary is not in --context folder", func() {
			oc.ImportJavaIS(commonVar.Project)
			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)

			newContext := helper.CreateNewContext()
			defer helper.DeleteDir(newContext)

			output := helper.CmdShouldFail("odo", append(args, "create", "--s2i", "java:8", "sb-jar-test", "--project",
				commonVar.Project, "--binary", filepath.Join(commonVar.Context, "sb.jar"), "--context", newContext)...)
			Expect(output).To(ContainSubstring("inside of the context directory"))
		})

		It("binary component is valid if path is relative and includes ../", func() {
			oc.ImportJavaIS(commonVar.Project)
			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)

			relativeContext := fmt.Sprintf("..%c%s", filepath.Separator, filepath.Base(commonVar.Context))
			fmt.Printf("relativeContext = %#v\n", relativeContext)

			if runtime.GOOS == "darwin" {
				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "java:8", "sb-jar-test", "--project",
					commonVar.Project, "--binary", filepath.Join("/private", commonVar.Context, "sb.jar"), "--context", relativeContext)...)
			} else {
				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "java:8", "sb-jar-test", "--project",
					commonVar.Project, "--binary", filepath.Join(commonVar.Context, "sb.jar"), "--context", relativeContext)...)
			}
			helper.ValidateLocalCmpExist(relativeContext, "Type,java:8", "Name,sb-jar-test")
		})

		It("should fail the create command as --git flag, which is specific to s2i component creation, is used without --s2i flag", func() {
			output := helper.CmdShouldFail("odo", "create", "nodejs", "cmp-git", "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")
			Expect(output).Should(ContainSubstring("flag --git, requires --s2i flag to be set, when deploying S2I (Source-to-Image) components"))
		})

		It("should fail the create command as --binary flag, which is specific to s2i component creation, is used without --s2i flag", func() {
			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)

			output := helper.CmdShouldFail("odo", "create", "java:8", "sb-jar-test", "--binary", filepath.Join(commonVar.Context, "sb.jar"), "--context", commonVar.Context)
			Expect(output).Should(ContainSubstring("flag --binary, requires --s2i flag to be set, when deploying S2I (Source-to-Image) components"))
		})
	})

	Context("Test odo push with --source and --config flags", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})
		Context("Using project flag(--project) and current directory", func() {
			It("create local nodejs component and push source and code separately", func() {
				appName := "nodejs-push-test"
				cmpName := "nodejs"
				helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)

				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project)...)
				helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)

				// component doesn't exist yet so attempt to only push source should fail
				helper.CmdShouldFail("odo", append(args, "push", "--source")...)

				// Push only config and see that the component is created but without any source copied
				helper.CmdShouldPass("odo", append(args, "push", "--config")...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)

				// Push only source and see that the component is updated with source code
				helper.CmdShouldPass("odo", append(args, "push", "--source")...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)
				remoteCmdExecPass := oc.CheckCmdOpInRemoteCmpPod(
					cmpName,
					appName,
					commonVar.Project,
					[]string{"sh", "-c", "ls -la $ODO_S2I_DEPLOYMENT_DIR/package.json"},
					func(cmdOp string, err error) bool {
						return err == nil
					},
				)
				Expect(remoteCmdExecPass).To(Equal(true))
			})

			It("create local nodejs component and push source and code at once", func() {
				appName := "nodejs-push-test"
				cmpName := "nodejs-push-atonce"
				helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)

				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project)...)
				helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)
				// Push only config and see that the component is created but without any source copied
				helper.CmdShouldPass("odo", append(args, "push")...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)
				remoteCmdExecPass := oc.CheckCmdOpInRemoteCmpPod(
					cmpName,
					appName,
					commonVar.Project,
					[]string{"sh", "-c", "ls -la $ODO_S2I_DEPLOYMENT_DIR/package.json"},
					func(cmdOp string, err error) bool {
						return err == nil
					},
				)
				Expect(remoteCmdExecPass).To(Equal(true))
			})

		})

		Context("when --context is used", func() {
			// don't need to switch to any dir here, as this test should use --context flag
			It("create local nodejs component and push source and code separately", func() {
				appName := "nodejs-push-context-test"
				cmpName := "nodejs"
				helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)

				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--context", commonVar.Context, "--app", appName, "--project", commonVar.Project)...)
				helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)

				// component doesn't exist yet so attempt to only push source should fail
				helper.CmdShouldFail("odo", append(args, "push", "--source", "--context", commonVar.Context)...)

				// Push only config and see that the component is created but without any source copied
				helper.CmdShouldPass("odo", append(args, "push", "--config", "--context", commonVar.Context)...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)

				// Push only source and see that the component is updated with source code
				helper.CmdShouldPass("odo", append(args, "push", "--source", "--context", commonVar.Context)...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)
				remoteCmdExecPass := oc.CheckCmdOpInRemoteCmpPod(
					cmpName,
					appName,
					commonVar.Project,
					[]string{"sh", "-c", "ls -la $ODO_S2I_DEPLOYMENT_DIR/package.json"},
					func(cmdOp string, err error) bool {
						return err == nil
					},
				)
				Expect(remoteCmdExecPass).To(Equal(true))
			})

			It("create local nodejs component and push source and code at once", func() {
				appName := "nodejs-push-context-test"
				cmpName := "nodejs-push-atonce"
				helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)

				helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--context", commonVar.Context, "--project", commonVar.Project)...)
				helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)

				// Push both config and source
				helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)
				oc.VerifyCmpExists(cmpName, appName, commonVar.Project)
				remoteCmdExecPass := oc.CheckCmdOpInRemoteCmpPod(
					cmpName,
					appName,
					commonVar.Project,
					[]string{"sh", "-c", "ls -la $ODO_S2I_DEPLOYMENT_DIR/package.json"},
					func(cmdOp string, err error) bool {
						return err == nil
					},
				)
				Expect(remoteCmdExecPass).To(Equal(true))
			})
		})
	})

	Context("Test odo push with --now flag during creation", func() {
		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})
		It("should successfully create config and push code in one create command with --now", func() {
			appName := "nodejs-create-now-test"
			cmpName := "nodejs-push-atonce"
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--now")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)

			oc.VerifyCmpExists(cmpName, appName, commonVar.Project)
			remoteCmdExecPass := oc.CheckCmdOpInRemoteCmpPod(
				cmpName,
				appName,
				commonVar.Project,
				[]string{"sh", "-c", "ls -la $ODO_S2I_DEPLOYMENT_DIR/package.json"},
				func(cmdOp string, err error) bool {
					return err == nil
				},
			)
			Expect(remoteCmdExecPass).To(Equal(true))
		})
	})

	Context("when component is in the current directory and --project flag is used", func() {

		appName := "app"
		componentName := "my-component"

		JustBeforeEach(func() {
			helper.Chdir(commonVar.Context)
		})

		It("create local nodejs component twice and fail", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "--env", "key=value,key1=value1")...)
			output := helper.CmdShouldFail("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "--env", "key=value,key1=value1")...)
			Expect(output).To(ContainSubstring("this directory already contains a component"))
		})

		It("creates and pushes local nodejs component and then deletes --all", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", componentName, "--app", appName, "--project", commonVar.Project, "--env", "key=value,key1=value1")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+componentName, "Application,"+appName)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", append(args, "delete", "--context", commonVar.Context, "-f", "--all")...)
			componentList := helper.CmdShouldPass("odo", append(args, "list", "--app", appName, "--project", commonVar.Project)...)
			Expect(componentList).NotTo(ContainSubstring(componentName))
			files := helper.ListFilesInDir(commonVar.Context)
			Expect(files).NotTo(ContainElement(".odo"))
		})

		It("creates a local python component, pushes it and then deletes it using --all flag", func() {
			helper.CopyExample(filepath.Join("source", "python"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "python", componentName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,python", "Name,"+componentName, "Application,"+appName)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", append(args, "delete", "--context", commonVar.Context, "-f")...)
			helper.CmdShouldPass("odo", append(args, "delete", "--all", "-f", "--context", commonVar.Context)...)
			componentList := helper.CmdShouldPass("odo", append(args, "list", "--app", appName, "--project", commonVar.Project)...)
			Expect(componentList).NotTo(ContainSubstring(componentName))
			files := helper.ListFilesInDir(commonVar.Context)
			Expect(files).NotTo(ContainElement(".odo"))
		})

		It("creates a local python component, pushes it and then deletes it using --all flag in local directory", func() {
			helper.CopyExample(filepath.Join("source", "python"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "python", componentName, "--app", appName, "--project", commonVar.Project)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,python", "Name,"+componentName, "Application,"+appName)
			helper.CmdShouldPass("odo", append(args, "push")...)
			helper.CmdShouldPass("odo", append(args, "delete", "--all", "-f")...)
			componentList := helper.CmdShouldPass("odo", append(args, "list", "--app", appName, "--project", commonVar.Project)...)
			Expect(componentList).NotTo(ContainSubstring(componentName))
			files := helper.ListFilesInDir(commonVar.Context)
			Expect(files).NotTo(ContainElement(".odo"))
		})

		It("creates a local python component and check for unsupported warning", func() {
			helper.CopyExample(filepath.Join("source", "python"), commonVar.Context)
			output := helper.CmdShouldPass("odo", append(args, "create", "--s2i", "python", componentName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			Expect(output).To(ContainSubstring("Warning: python is not fully supported by odo, and it is not guaranteed to work"))
		})

		It("creates a local nodejs component and check unsupported warning hasn't occurred", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			output := helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs:latest", componentName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			Expect(output).NotTo(ContainSubstring("Warning"))
		})

		It("creates a local java component and check unsupported warning hasn't occurred", func() {
			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)
			output := helper.CmdShouldPass("odo", append(args, "create", "--s2i", "java:latest", componentName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			Expect(output).NotTo(ContainSubstring("Warning"))
		})
	})

	Context("odo component updating", func() {

		It("should be able to create a git component and update it from local to git", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--context", commonVar.Context, "--app", "testing")...)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			helper.CmdShouldPass("odo", "update", "--git", "https://github.com/openshift/nodejs-ex.git", "--context", commonVar.Context)
			// check the source location and type in the deployment config
			getSourceLocation := oc.SourceLocationDC("cmp-git", "testing", commonVar.Project)
			Expect(getSourceLocation).To(ContainSubstring("https://github.com/openshift/nodejs-ex"))
			getSourceType := oc.SourceTypeDC("cmp-git", "testing", commonVar.Project)
			Expect(getSourceType).To(ContainSubstring("git"))
		})

		It("should be able to update a component from git to local", func() {
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "cmp-git", "--project", commonVar.Project, "--git", "https://github.com/openshift/nodejs-ex", "--context", commonVar.Context, "--app", "testing")...)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			// update the component config according to the git component
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)

			helper.CmdShouldPass("odo", "update", "--local", "./", "--context", commonVar.Context)

			// check the source location and type in the deployment config
			getSourceLocation := oc.SourceLocationDC("cmp-git", "testing", commonVar.Project)
			Expect(getSourceLocation).To(ContainSubstring(""))
			getSourceType := oc.SourceTypeDC("cmp-git", "testing", commonVar.Project)
			Expect(getSourceType).To(ContainSubstring("local"))
		})
	})

	Context("odo component delete, list and describe", func() {
		appName := "app"
		cmpName := "nodejs"

		It("should pass inside a odo directory without component name as parameter", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", "url", "create", "example", "--context", commonVar.Context)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName, "URL,0,Name,example")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			// changing directory to the context directory
			helper.Chdir(commonVar.Context)
			cmpListOutput := helper.CmdShouldPass("odo", append(args, "list")...)
			Expect(cmpListOutput).To(ContainSubstring(cmpName))
			cmpDescribe := helper.CmdShouldPass("odo", append(args, "describe")...)
			helper.MatchAllInOutput(cmpDescribe, []string{cmpName, "nodejs"})

			url := helper.DetermineRouteURL(commonVar.Context)
			Expect(cmpDescribe).To(ContainSubstring(url))

			helper.CmdShouldPass("odo", append(args, "delete", "-f")...)
		})

		It("should fail outside a odo directory without component name as parameter", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			// commands should fail as the component name is missing
			helper.CmdShouldFail("odo", append(args, "describe", "--app", appName, "--project", commonVar.Project)...)
			helper.CmdShouldFail("odo", append(args, "delete", "-f", "--app", appName, "--project", commonVar.Project)...)
		})

		It("should pass outside a odo directory with component name as parameter", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			cmpListOutput := helper.CmdShouldPass("odo", append(args, "list", "--app", appName, "--project", commonVar.Project)...)
			Expect(cmpListOutput).To(ContainSubstring(cmpName))

			actualDesCompJSON := helper.CmdShouldPass("odo", append(args, "describe", cmpName, "--app", appName, "--project", commonVar.Project, "-o", "json")...)

			desiredDesCompJSON := fmt.Sprintf(`{"kind":"Component","apiVersion":"odo.dev/v1alpha1","metadata":{"name":"nodejs","namespace":"%s","creationTimestamp":null},"spec":{"app":"app","type":"nodejs","sourceType": "local", "urls": {"kind": "List", "apiVersion": "odo.dev/v1alpha1", "metadata": {}, "items": null}, "storages": {"kind": "List", "apiVersion": "odo.dev/v1alpha1", "metadata": {}, "items": null}, "env":[{"name":"DEBUG_PORT","value":"5858"}]},"status":{"state":"Pushed"}}`, commonVar.Project)
			Expect(desiredDesCompJSON).Should(MatchJSON(actualDesCompJSON))

			helper.CmdShouldPass("odo", append(args, "delete", cmpName, "--app", appName, "--project", commonVar.Project, "-f")...)
		})
	})

	Context("when running odo push multiple times, check for existence of environment variables", func() {
		It("should should retain the same environment variable on multiple push", func() {
			componentName := helper.RandString(6)
			appName := helper.RandString(6)
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", componentName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			helper.Chdir(commonVar.Context)
			helper.CmdShouldPass("odo", "config", "set", "--env", "FOO=BAR")
			helper.CmdShouldPass("odo", append(args, "push")...)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+componentName, "Application,"+appName, "Ports,[8080/TCP]", "Envs,0,Name,FOO")
			ports := oc.GetDcPorts(componentName, appName, commonVar.Project)
			Expect(ports).To(ContainSubstring("8080"))
			dcName := oc.GetDcName(componentName, commonVar.Project)
			stdOut := helper.CmdShouldPass("oc", "get", "dc/"+dcName, "-n", commonVar.Project, "-o", "go-template={{ .spec.template.spec }}{{.env}}")
			Expect(stdOut).To(ContainSubstring("FOO"))

			helper.CmdShouldPass("odo", append(args, "push")...)
			stdOut = oc.DescribeDc(dcName, commonVar.Project)
			Expect(stdOut).To(ContainSubstring("FOO"))
		})
	})

	Context("Creating component with numeric named context", func() {
		var contextNumeric string
		JustBeforeEach(func() {
			var err error
			ts := time.Now().UnixNano()
			contextNumeric, err = ioutil.TempDir("", fmt.Sprint(ts))
			Expect(err).ToNot(HaveOccurred())
		})
		JustAfterEach(func() {
			helper.DeleteDir(contextNumeric)
		})

		It("should create default named component in a directory with numeric name", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), contextNumeric)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", "--project", commonVar.Project, "--context", contextNumeric, "--app", "testing")...)
			helper.ValidateLocalCmpExist(contextNumeric, "Type,nodejs", "Application,testing")
			helper.CmdShouldPass("odo", append(args, "push", "--context", contextNumeric, "-v4")...)
		})
	})

	Context("Creating component using symlink", func() {
		var symLinkPath string

		JustBeforeEach(func() {
			if runtime.GOOS == "windows" {
				Skip("Skipping test because for symlink creation on platform like Windows, go library needs elevated privileges.")
			}
			// create a symlink
			symLinkName := helper.RandString(10)
			helper.CreateSymLink(commonVar.Context, filepath.Join(filepath.Dir(commonVar.Context), symLinkName))
			symLinkPath = filepath.Join(filepath.Dir(commonVar.Context), symLinkName)
		})
		JustAfterEach(func() {
			// remove the symlink
			err := os.Remove(symLinkPath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should be able to deploy a spring boot uberjar file using symlinks in all odo commands", func() {
			oc.ImportJavaIS(commonVar.Project)

			helper.CopyExample(filepath.Join("binary", "java", "openjdk"), commonVar.Context)

			// create the component using symlink
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "java:8", "sb-jar-test", "--project",
				commonVar.Project, "--binary", filepath.Join(symLinkPath, "sb.jar"), "--context", symLinkPath)...)

			// Create a URL and push without using the symlink
			helper.CmdShouldPass("odo", "url", "create", "uberjaropenjdk", "--port", "8080", "--context", symLinkPath)
			helper.ValidateLocalCmpExist(symLinkPath, "Type,java:8", "Name,sb-jar-test", "Application,app", "URL,0,Name,uberjaropenjdk")
			helper.CmdShouldPass("odo", append(args, "push", "--context", symLinkPath)...)
			routeURL := helper.DetermineRouteURL(symLinkPath)

			// Ping said URL
			helper.HttpWaitFor(routeURL, "HTTP Booster", 300, 1)

			// Delete the component
			helper.CmdShouldPass("odo", append(args, "delete", "sb-jar-test", "-f", "--context", symLinkPath)...)
		})

		It("Should be able to deploy a wildfly war file using symlinks in some odo commands", func() {
			helper.CopyExample(filepath.Join("binary", "java", "wildfly"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "wildfly", "javaee-war-test", "--project",
				commonVar.Project, "--binary", filepath.Join(symLinkPath, "ROOT.war"), "--context", symLinkPath)...)

			// Create a URL
			helper.CmdShouldPass("odo", "url", "create", "warfile", "--port", "8080", "--context", commonVar.Context)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,wildfly", "Name,javaee-war-test", "Application,app", "URL,0,Name,warfile")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)
			routeURL := helper.DetermineRouteURL(commonVar.Context)

			// Ping said URL
			helper.HttpWaitFor(routeURL, "Sample", 90, 1)

			// Delete the component
			helper.CmdShouldPass("odo", append(args, "delete", "javaee-war-test", "-f", "--context", commonVar.Context)...)
		})
	})

	Context("odo component delete should clean owned resources", func() {
		appName := helper.RandString(5)
		cmpName := helper.RandString(5)

		It("should delete the component and the owned resources", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", "url", "create", "example-1", "--context", commonVar.Context)

			helper.CmdShouldPass("odo", "storage", "create", "storage-1", "--size", "1Gi", "--path", "/data1", "--context", commonVar.Context)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName, "URL,0,Name,example-1")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			helper.CmdShouldPass("odo", "url", "create", "example-2", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "storage", "create", "storage-2", "--size", "1Gi", "--path", "/data2", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			helper.CmdShouldPass("odo", append(args, "delete", "-f", "--context", commonVar.Context)...)

			oc.WaitAndCheckForExistence("routes", commonVar.Project, 1)
			oc.WaitAndCheckForExistence("dc", commonVar.Project, 1)
			oc.WaitAndCheckForExistence("pvc", commonVar.Project, 1)
			oc.WaitAndCheckForExistence("bc", commonVar.Project, 1)
			oc.WaitAndCheckForExistence("is", commonVar.Project, 1)
			oc.WaitAndCheckForExistence("service", commonVar.Project, 1)
		})

		It("should delete the component and the owned resources with wait flag", func() {
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "create", "--s2i", "nodejs", cmpName, "--app", appName, "--project", commonVar.Project, "--context", commonVar.Context)...)
			helper.CmdShouldPass("odo", "url", "create", "example-1", "--context", commonVar.Context)

			helper.CmdShouldPass("odo", "storage", "create", "storage-1", "--size", "1Gi", "--path", "/data1", "--context", commonVar.Context)
			helper.ValidateLocalCmpExist(commonVar.Context, "Type,nodejs", "Name,"+cmpName, "Application,"+appName, "URL,0,Name,example-1")
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			helper.CmdShouldPass("odo", "url", "create", "example-2", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "storage", "create", "storage-2", "--size", "1Gi", "--path", "/data2", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", append(args, "push", "--context", commonVar.Context)...)

			// delete with --wait flag
			helper.CmdShouldPass("odo", append(args, "delete", "-f", "-w", "--context", commonVar.Context)...)

			oc.VerifyResourceDeleted("routes", "example", commonVar.Project)
			oc.VerifyResourceDeleted("service", cmpName, commonVar.Project)
			// verify s2i pvc is delete
			oc.VerifyResourceDeleted("pvc", "s2idata", commonVar.Project)
			oc.VerifyResourceDeleted("pvc", "storage-1", commonVar.Project)
			oc.VerifyResourceDeleted("pvc", "storage-2", commonVar.Project)
			oc.VerifyResourceDeleted("dc", cmpName, commonVar.Project)
		})

	})

	Context("convert s2i to devfile", func() {

		JustBeforeEach(func() {
			os.Setenv("ODO_EXPERIMENTAL", "true")
		})

		JustAfterEach(func() {
			os.Unsetenv("ODO_EXPERIMENTAL")
		})

		It("should convert s2i component to devfile component successfully", func() {
			cmpName := "mynodejs"
			appName := "app"
			urlName := "url1"
			storageName := "storage1"

			// create a s2i component
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CmdShouldPass("odo", "component", "create", "--s2i", "nodejs", cmpName, "--project", commonVar.Project, "--context", commonVar.Context, "--app", appName, "--s2i")
			helper.CmdShouldPass("odo", "url", "create", urlName, "--port", "8080", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "storage", "create", storageName, "--path", "/data1", "--size", "1Gi", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)

			// convert it to devfile
			helper.CmdShouldPass("odo", "utils", "convert-to-devfile", "--context", commonVar.Context)
			helper.CmdShouldPass("odo", "push", "--context", commonVar.Context)

			// check the status of devfile component
			stdout := helper.CmdShouldPass("odo", "list", "--context", commonVar.Context)
			helper.MatchAllInOutput(stdout, []string{cmpName, "Devfile Components", "Pushed"})

			// delete the s2i component
			helper.CmdShouldPass("odo", "delete", "--s2i", "-a", "-f", "--context", commonVar.Context)

			// verify the url
			stdout = helper.CmdShouldPass("odo", "url", "list", "--context", commonVar.Context)

			helper.MatchAllInOutput(stdout, []string{urlName, "Pushed", "false", "route"})
			//verify storage
			stdout = helper.CmdShouldPass("odo", "storage", "list", "--context", commonVar.Context)
			helper.MatchAllInOutput(stdout, []string{storageName, "Pushed"})

		})

	})

}
