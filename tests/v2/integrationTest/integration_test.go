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

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/v2/pkg/attributes"
	commonUtils "github.com/devfile/api/v2/test/v200/utils/common"
	"github.com/devfile/library/pkg/devfile/parser"
	"github.com/devfile/library/pkg/testingutil"
	libraryUtils "github.com/devfile/library/tests/v2/utils/library"
	kubev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getInvalidNodeJSDevfileList() []string {
	return []string{
		"devfile-with-invalid-events.yaml",
		"devfile-with-invalid-volmount.yaml",
		"devfile-with-multiple-defaults.yaml",
		"devfile-with-no-default.yaml",
		"devfile-with-preStart.yaml",
		"devfileCompositeInvalidComponent.yaml",
		"devfileCompositeNonExistent.yaml",
		"devfileIndirectNesting.yaml",
	}
}

func IsValidNodeJSDevFile(fileName string) bool {

	fileNames := getInvalidNodeJSDevfileList()
	for _, item := range fileNames {
		if item == fileName {
			return false
		}
	}
	return true
}

func getTestDevfileList(subdir string, testValid bool) ([]string, error) {

	fileList := make([]string, 0)
	subdirPath := "../../examples/source/devfiles/" + subdir + "/"
	files, err := ioutil.ReadDir(subdirPath)

	if err != nil {
		commonUtils.LogErrorMessage(fmt.Sprintf("Error in getting file list from the directory: %s : %v", subdirPath, err))
	}

	for _, file := range files {
		if !file.IsDir() {
			r, err := regexp.MatchString("^devfile.+yaml$", file.Name())

			if err == nil && r {
				if subdir == "nodejs" {
					if testValid {
						if IsValidNodeJSDevFile(file.Name()) {
							fileList = append(fileList, file.Name())
						}
					} else {
						if !IsValidNodeJSDevFile(file.Name()) {
							fileList = append(fileList, file.Name())
						}
					}

				} else {
					fileList = append(fileList, file.Name())
				}
			}
		}
	}
	return fileList, err
}

func testDevfile(t *testing.T, subdir string, testValid bool) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = false
	testContent.EditContent = false

	fileNames, _ := getTestDevfileList(subdir, testValid)

	libraryUtils.CopyTestDevfile(t, subdir, fileNames)

	for _, fileName := range fileNames {
		testContent.FileName = fileName
		if testValid {
			libraryUtils.RunStaticTest(testContent, t)
		} else {
			libraryUtils.RunStaticTestToFail(testContent, t)
		}
	}
}

func Test_Valid_NodeJS_Devfiles(t *testing.T) {
	testDevfile(t, "nodejs", true)
}

func Test_Invalid_NodeJS_Devfiles(t *testing.T) {
	testDevfile(t, "nodejs", false)
}

func Test_Valid_OpenLiberty_Devfiles(t *testing.T) {
	testDevfile(t, "java-openliberty", true)
}

func Test_Valid_Python_Devfiles(t *testing.T) {
	testDevfile(t, "python", true)
}

func Test_Valid_Springboot_Devfiles(t *testing.T) {
	testDevfile(t, "springboot", true)
}

func Test_Parent_Local_URI(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = true
	testContent.EditContent = false
	testContent.FileName = "Test_Parent_LocalURI.yaml"
	//copy the parent and main devfile from devfiles/samples
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName, "Parent.yaml"})
	libraryUtils.RunStaticTest(testContent, t)
	libraryUtils.RunMultiThreadedStaticTest(testContent, t)
}

func Test_v200_Devfile(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = false
	testContent.EditContent = false
	testContent.FileName = "Test_200.yaml"
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName})
	libraryUtils.RunStaticTest(testContent, t)
	libraryUtils.RunMultiThreadedStaticTest(testContent, t)
}

func Test_v210_Devfile(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = false
	testContent.EditContent = false
	testContent.FileName = "Test_210.yaml"
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName})
	libraryUtils.RunStaticTest(testContent, t)
	libraryUtils.RunMultiThreadedStaticTest(testContent, t)
}

func Test_v220_Devfile(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = false
	testContent.EditContent = false
	testContent.FileName = "Test_220.yaml"
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName})
	libraryUtils.RunStaticTest(testContent, t)

}

//Create kube client and context and set as ParserArgs for Parent Kubernetes reference test.  Corresponding main devfile is ../devfile/samples/TestParent_KubeCRD.yaml
func setClientAndContextParserArgs() *parser.ParserArgs {
	isTrue := true
	name := "testkubeparent1"
	parentSpec := schema.DevWorkspaceTemplateSpec{
		DevWorkspaceTemplateSpecContent: schema.DevWorkspaceTemplateSpecContent{
			Commands: []schema.Command{
				{
					Id: "applycommand",
					CommandUnion: schema.CommandUnion{
						Apply: &schema.ApplyCommand{
							Component: "devbuild",
							LabeledCommand: schema.LabeledCommand{
								Label: "testcontainerparent",
								BaseCommand: schema.BaseCommand{
									Group: &schema.CommandGroup{
										Kind:      schema.TestCommandGroupKind,
										IsDefault: &isTrue,
									},
								},
							},
						},
					},
				},
			},
			Components: []schema.Component{
				{
					Name: "devbuild",
					ComponentUnion: schema.ComponentUnion{
						Container: &schema.ContainerComponent{
							Container: schema.Container{
								Image: "quay.io/nodejs-12",
							},
						},
					},
				},
			},
			Projects: []schema.Project{
				{
					Name: "parentproject",
					ProjectSource: schema.ProjectSource{
						Git: &schema.GitProjectSource{
							GitLikeProjectSource: schema.GitLikeProjectSource{
								CheckoutFrom: &schema.CheckoutFrom{
									Revision: "master",
									Remote:   "origin",
								},
								Remotes: map[string]string{"origin": "https://github.com/spring-projects/spring-petclinic.git"},
							},
						},
					},
				},
				{
					Name: "parentproject2",
					ProjectSource: schema.ProjectSource{
						Zip: &schema.ZipProjectSource{
							Location: "https://github.com/spring-projects/spring-petclinic.zip",
						},
					},
				},
			},
			StarterProjects: []schema.StarterProject{
				{
					Name: "parentstarterproject",
					ProjectSource: schema.ProjectSource{
						Git: &schema.GitProjectSource{
							GitLikeProjectSource: schema.GitLikeProjectSource{
								CheckoutFrom: &schema.CheckoutFrom{
									Revision: "main",
									Remote:   "origin",
								},
								Remotes: map[string]string{"origin": "https://github.com/spring-projects/spring-petclinic.git"},
							},
						},
					},
				},
			},
			Attributes: attributes.Attributes{}.FromStringMap(map[string]string{"category": "parentDevfile", "title": "This is a parent devfile"}),
			Variables:  map[string]string{"version": "2.0.0", "tag": "parent"},
		},
	}
	testK8sClient := &testingutil.FakeK8sClient{
		DevWorkspaceResources: map[string]schema.DevWorkspaceTemplate{
			name: {
				TypeMeta: kubev1.TypeMeta{
					APIVersion: "2.1.0",
				},
				Spec: parentSpec,
			},
		},
	}
	parserArgs := parser.ParserArgs{}
	parserArgs.K8sClient = testK8sClient
	parserArgs.Context = context.Background()
	return &parserArgs
}

func Test_Parent_KubeCRD(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = true
	testContent.EditContent = false
	testContent.FileName = "Test_Parent_KubeCRD.yaml"
	parserArgs := setClientAndContextParserArgs()
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName})
	libraryUtils.SetParserArgs(*parserArgs)
	libraryUtils.RunStaticTest(testContent, t)
	libraryUtils.RunMultiThreadedStaticTest(testContent, t)
}

func Test_Parent_RegistryURL(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = true
	testContent.EditContent = false
	testContent.FileName = "Test_Parent_RegistryURL.yaml"
	libraryUtils.CopyDevfileSamples(t, []string{testContent.FileName})
	libraryUtils.RunStaticTest(testContent, t)
	libraryUtils.RunMultiThreadedStaticTest(testContent, t)
}

func Test_Everything(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = commonUtils.CommandTypes
	testContent.ComponentTypes = commonUtils.ComponentTypes
	testContent.ProjectTypes = commonUtils.ProjectSourceTypes
	testContent.StarterProjectTypes = commonUtils.ProjectSourceTypes
	testContent.AddEvents = true
	testContent.AddMetaData = true
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	libraryUtils.RunTest(testContent, t)
	libraryUtils.RunMultiThreadTest(testContent, t)
}

func Test_EverythingEdit(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = commonUtils.CommandTypes
	testContent.ComponentTypes = commonUtils.ComponentTypes
	testContent.ProjectTypes = commonUtils.ProjectSourceTypes
	testContent.StarterProjectTypes = commonUtils.ProjectSourceTypes
	testContent.AddEvents = true
	testContent.AddMetaData = true
	testContent.EditContent = true
	testContent.FileName = commonUtils.GetDevFileName()
	libraryUtils.RunTest(testContent, t)
	libraryUtils.RunMultiThreadTest(testContent, t)
}
