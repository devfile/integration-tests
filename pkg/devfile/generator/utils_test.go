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

package generator

import (
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/api/v2/pkg/attributes"
	"github.com/devfile/library/pkg/devfile/parser"
	"github.com/devfile/library/pkg/devfile/parser/data"
	"github.com/devfile/library/pkg/devfile/parser/data/v2/common"
	"github.com/devfile/library/pkg/testingutil"
	"github.com/golang/mock/gomock"
	buildv1 "github.com/openshift/api/build/v1"

	v1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var isTrue bool = true

func TestConvertEnvs(t *testing.T) {
	envVarsNames := []string{"test", "sample-var", "myvar"}
	envVarsValues := []string{"value1", "value2", "value3"}
	tests := []struct {
		name    string
		envVars []v1.EnvVar
		want    []corev1.EnvVar
	}{
		{
			name: "One env var",
			envVars: []v1.EnvVar{
				{
					Name:  envVarsNames[0],
					Value: envVarsValues[0],
				},
			},
			want: []corev1.EnvVar{
				{
					Name:  envVarsNames[0],
					Value: envVarsValues[0],
				},
			},
		},
		{
			name: "Multiple env vars",
			envVars: []v1.EnvVar{
				{
					Name:  envVarsNames[0],
					Value: envVarsValues[0],
				},
				{
					Name:  envVarsNames[1],
					Value: envVarsValues[1],
				},
				{
					Name:  envVarsNames[2],
					Value: envVarsValues[2],
				},
			},
			want: []corev1.EnvVar{
				{
					Name:  envVarsNames[0],
					Value: envVarsValues[0],
				},
				{
					Name:  envVarsNames[1],
					Value: envVarsValues[1],
				},
				{
					Name:  envVarsNames[2],
					Value: envVarsValues[2],
				},
			},
		},
		{
			name:    "No env vars",
			envVars: []v1.EnvVar{},
			want:    []corev1.EnvVar{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVars := convertEnvs(tt.envVars)
			if !reflect.DeepEqual(tt.want, envVars) {
				t.Errorf("TestConvertEnvs() error: expected %v, wanted %v", envVars, tt.want)
			}
		})
	}
}

func TestConvertPorts(t *testing.T) {
	endpointsNames := []string{"endpoint1", "endpoint2", "a-very-long-port-name-before-endpoint-length-limit-8080"}
	endpointsPorts := []int{8080, 9090}
	tests := []struct {
		name      string
		endpoints []v1.Endpoint
		want      []corev1.ContainerPort
	}{
		{
			name: "One Endpoint",
			endpoints: []v1.Endpoint{
				{
					Name:       endpointsNames[0],
					TargetPort: endpointsPorts[0],
				},
			},
			want: []corev1.ContainerPort{
				{
					Name:          endpointsNames[0],
					ContainerPort: int32(endpointsPorts[0]),
					Protocol:      "TCP",
				},
			},
		},
		{
			name: "One Endpoint with >15 chars length",
			endpoints: []v1.Endpoint{
				{
					Name:       endpointsNames[2],
					TargetPort: endpointsPorts[0],
				},
			},
			want: []corev1.ContainerPort{
				{
					Name:          "port-8080",
					ContainerPort: int32(endpointsPorts[0]),
					Protocol:      "TCP",
				},
			},
		},
		{
			name: "Multiple endpoints",
			endpoints: []v1.Endpoint{
				{
					Name:       endpointsNames[0],
					TargetPort: endpointsPorts[0],
				},
				{
					Name:       endpointsNames[1],
					TargetPort: endpointsPorts[1],
				},
			},
			want: []corev1.ContainerPort{
				{
					Name:          endpointsNames[0],
					ContainerPort: int32(endpointsPorts[0]),
					Protocol:      "TCP",
				},
				{
					Name:          endpointsNames[1],
					ContainerPort: int32(endpointsPorts[1]),
					Protocol:      "TCP",
				},
			},
		},
		{
			name:      "No endpoints",
			endpoints: []v1.Endpoint{},
			want:      []corev1.ContainerPort{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ports := convertPorts(tt.endpoints)
			if !reflect.DeepEqual(tt.want, ports) {
				t.Errorf("TestConvertPorts() error: expected %v, wanted %v", ports, tt.want)
			}
		})
	}
}

func TestGetResourceReqs(t *testing.T) {
	memoryLimit := "1024Mi"
	memoryRequest := "1Gi"
	cpuRequest := "1m"
	cpuLimit := "1m"

	memoryLimitQuantity, err := resource.ParseQuantity(memoryLimit)
	memoryRequestQuantity, err := resource.ParseQuantity(memoryRequest)
	cpuRequestQuantity, err := resource.ParseQuantity(cpuRequest)
	cpuLimitQuantity, err := resource.ParseQuantity(cpuLimit)
	if err != nil {
		t.Errorf("TestGetResourceReqs() unexpected error: %v", err)
	}
	tests := []struct {
		name      string
		component v1.Component
		want      corev1.ResourceRequirements
		wantErr   []string
	}{
		{
			name: "generate resource limit",
			component: v1.Component{
				Name: "testcomponent",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							MemoryLimit:   memoryLimit,
							MemoryRequest: memoryRequest,
							CpuRequest:    cpuRequest,
							CpuLimit:      cpuLimit,
						},
					},
				},
			},
			want: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: memoryLimitQuantity,
					corev1.ResourceCPU:    cpuLimitQuantity,
				},
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: memoryRequestQuantity,
					corev1.ResourceCPU:    cpuRequestQuantity,
				},
			},
		},
		{
			name:      "Empty Component",
			component: v1.Component{},
			want:      corev1.ResourceRequirements{},
		},
		{
			name: "Valid container, but empty memoryLimit",
			component: v1.Component{
				Name: "testcomponent",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							Image: "testimage",
						},
					},
				},
			},
			want: corev1.ResourceRequirements{},
		},
		{
			name: "test error case",
			component: v1.Component{
				Name: "testcomponent",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							MemoryLimit:   "invalid",
							MemoryRequest: "invalid",
							CpuRequest:    "invalid",
							CpuLimit:      "invalid",
						},
					},
				},
			},
			wantErr: []string{
				"error parsing memoryLimit requirement.*",
				"error parsing cpuLimit requirement.*",
				"error parsing memoryRequest requirement.*",
				"error parsing cpuRequest requirement.*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := getResourceReqs(tt.component)
			if merr, ok := err.(*multierror.Error); ok && tt.wantErr != nil {
				assert.Equal(t, len(tt.wantErr), len(merr.Errors), "Error list length should match")
				for i := 0; i < len(merr.Errors); i++ {
					assert.Regexp(t, tt.wantErr[i], merr.Errors[i].Error(), "Error message should match")
				}
			} else if !reflect.DeepEqual(tt.want, req) {
				assert.Equal(t, tt.want, req, "TestGetResourceReqs(): The two values should be the same.")
			}
		})
	}
}

func TestAddSyncRootFolder(t *testing.T) {

	tests := []struct {
		name               string
		sourceMapping      string
		wantSyncRootFolder string
	}{
		{
			name:               "Valid Source Mapping",
			sourceMapping:      "/mypath",
			wantSyncRootFolder: "/mypath",
		},
		{
			name:               "No Source Mapping",
			sourceMapping:      "",
			wantSyncRootFolder: DevfileSourceVolumeMount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := testingutil.CreateFakeContainer("container1")

			syncRootFolder := addSyncRootFolder(&container, tt.sourceMapping)

			if syncRootFolder != tt.wantSyncRootFolder {
				t.Errorf("TestAddSyncRootFolder() sync root folder error: expected %v got %v", tt.wantSyncRootFolder, syncRootFolder)
			}

			for _, env := range container.Env {
				if env.Name == EnvProjectsRoot && env.Value != tt.wantSyncRootFolder {
					t.Errorf("TestAddSyncRootFolder() PROJECT_ROOT error: expected %s, actual %s", tt.wantSyncRootFolder, env.Value)
				}
			}
		})
	}
}

func TestAddSyncFolder(t *testing.T) {
	projectNames := []string{"some-name", "another-name"}
	projectRepos := []string{"https://github.com/some/repo.git", "https://github.com/another/repo.git"}
	projectClonePath := "src/github.com/golang/example/"
	invalidClonePaths := []string{"/var", "../var", "pkg/../../var"}
	sourceVolumePath := "/projects/app"

	absoluteClonePathErr := "the clonePath .* in the devfile project .* must be a relative path"
	escapeClonePathErr := "the clonePath .* in the devfile project .* cannot escape the value defined by [$]PROJECTS_ROOT. Please avoid using \"..\" in clonePath"

	tests := []struct {
		name     string
		projects []v1.Project
		want     string
		wantErr  *string
	}{
		{
			name:     "No projects",
			projects: []v1.Project{},
			want:     sourceVolumePath,
		},
		{
			name: "One project",
			projects: []v1.Project{
				{
					Name: projectNames[0],
					ProjectSource: v1.ProjectSource{
						Git: &v1.GitProjectSource{
							GitLikeProjectSource: v1.GitLikeProjectSource{
								Remotes: map[string]string{"origin": projectRepos[0]},
							},
						},
					},
				},
			},
			want: filepath.ToSlash(filepath.Join(sourceVolumePath, projectNames[0])),
		},
		{
			name: "Multiple projects",
			projects: []v1.Project{
				{
					Name: projectNames[0],
					ProjectSource: v1.ProjectSource{
						Git: &v1.GitProjectSource{
							GitLikeProjectSource: v1.GitLikeProjectSource{
								Remotes: map[string]string{"origin": projectRepos[0]},
							},
						},
					},
				},
				{
					Name: projectNames[1],
					ProjectSource: v1.ProjectSource{
						Zip: &v1.ZipProjectSource{
							Location: projectRepos[1],
						},
					},
				},
			},
			want: filepath.ToSlash(filepath.Join(sourceVolumePath, projectNames[0])),
		},
		{
			name: "Clone path set",
			projects: []v1.Project{
				{
					ClonePath: projectClonePath,
					Name:      projectNames[0],
					ProjectSource: v1.ProjectSource{
						Zip: &v1.ZipProjectSource{
							Location: projectRepos[0],
						},
					},
				},
			},
			want: filepath.ToSlash(filepath.Join(sourceVolumePath, projectClonePath)),
		},
		{
			name: "Invalid clone path, set with absolute path",
			projects: []v1.Project{
				{
					ClonePath: invalidClonePaths[0],
					Name:      projectNames[0],
					ProjectSource: v1.ProjectSource{
						Git: &v1.GitProjectSource{
							GitLikeProjectSource: v1.GitLikeProjectSource{
								Remotes: map[string]string{"origin": projectRepos[0]},
							},
						},
					},
				},
			},
			want:    "",
			wantErr: &absoluteClonePathErr,
		},
		{
			name: "Invalid clone path, starts with ..",
			projects: []v1.Project{
				{
					ClonePath: invalidClonePaths[1],
					Name:      projectNames[0],
					ProjectSource: v1.ProjectSource{
						Git: &v1.GitProjectSource{
							GitLikeProjectSource: v1.GitLikeProjectSource{
								Remotes: map[string]string{"origin": projectRepos[0]},
							},
						},
					},
				},
			},
			want:    "",
			wantErr: &escapeClonePathErr,
		},
		{
			name: "Invalid clone path, contains ..",
			projects: []v1.Project{
				{
					ClonePath: invalidClonePaths[2],
					Name:      projectNames[0],
					ProjectSource: v1.ProjectSource{
						Zip: &v1.ZipProjectSource{
							Location: projectRepos[0],
						},
					},
				},
			},
			want:    "",
			wantErr: &escapeClonePathErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := testingutil.CreateFakeContainer("container1")

			err := addSyncFolder(&container, sourceVolumePath, tt.projects)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("TestAddSyncFolder() error: unexpected error %v, want %v", err, tt.wantErr)
			} else if err == nil {
				for _, env := range container.Env {
					if env.Name == EnvProjectsSrc && env.Value != tt.want {
						t.Errorf("TestAddSyncFolder() error: expected %s, actual %s", tt.want, env.Value)
					}
				}
			} else {
				assert.Regexp(t, *tt.wantErr, err.Error(), "TestAddSyncFolder(): Error message should match")
			}
		})
	}
}

func TestGetContainer(t *testing.T) {

	tests := []struct {
		name          string
		containerName string
		image         string
		isPrivileged  bool
		command       []string
		args          []string
		envVars       []corev1.EnvVar
		resourceReqs  corev1.ResourceRequirements
		ports         []corev1.ContainerPort
	}{
		{
			name:          "Empty container params",
			containerName: "",
			image:         "",
			isPrivileged:  false,
			command:       []string{},
			args:          []string{},
			envVars:       []corev1.EnvVar{},
			resourceReqs:  corev1.ResourceRequirements{},
			ports:         []corev1.ContainerPort{},
		},
		{
			name:          "Valid container params",
			containerName: "container1",
			image:         "quay.io/eclipse/che-java8-maven:nightly",
			isPrivileged:  true,
			command:       []string{"tail"},
			args:          []string{"-f", "/dev/null"},
			envVars: []corev1.EnvVar{
				{
					Name:  "test",
					Value: "123",
				},
			},
			resourceReqs: fakeResources,
			ports: []corev1.ContainerPort{
				{
					Name:          "port-9090",
					ContainerPort: 9090,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containerParams := containerParams{
				Name:         tt.containerName,
				Image:        tt.image,
				IsPrivileged: tt.isPrivileged,
				Command:      tt.command,
				Args:         tt.args,
				EnvVars:      tt.envVars,
				ResourceReqs: tt.resourceReqs,
				Ports:        tt.ports,
			}
			container := getContainer(containerParams)

			if container.Name != tt.containerName {
				t.Errorf("TestGetContainer() error: expected containerName %s, actual %s", tt.containerName, container.Name)
			}

			if container.Image != tt.image {
				t.Errorf("TestGetContainer() error: expected image %s, actual %s", tt.image, container.Image)
			}

			if tt.isPrivileged {
				if *container.SecurityContext.Privileged != tt.isPrivileged {
					t.Errorf("TestGetContainer() error: expected isPrivileged %t, actual %t", tt.isPrivileged, *container.SecurityContext.Privileged)
				}
			} else if tt.isPrivileged == false && container.SecurityContext != nil {
				t.Errorf("expected security context to be nil but it was defined")
			}

			if len(container.Command) != len(tt.command) {
				t.Errorf("TestGetContainer() error: expected command length %d, actual %d", len(tt.command), len(container.Command))
			} else {
				for i := range container.Command {
					if container.Command[i] != tt.command[i] {
						t.Errorf("TestGetContainer() error: expected command %s, actual %s", tt.command[i], container.Command[i])
					}
				}
			}

			if len(container.Args) != len(tt.args) {
				t.Errorf("TestGetContainer() error: expected container args length %d, actual %d", len(tt.args), len(container.Args))
			} else {
				for i := range container.Args {
					if container.Args[i] != tt.args[i] {
						t.Errorf("TestGetContainer() error: expected container args %s, actual %s", tt.args[i], container.Args[i])
					}
				}
			}

			if len(container.Env) != len(tt.envVars) {
				t.Errorf("TestGetContainer() error: expected container env length %d, actual %d", len(tt.envVars), len(container.Env))
			} else {
				for i := range container.Env {
					if container.Env[i].Name != tt.envVars[i].Name {
						t.Errorf("TestGetContainer() error: expected env name %s, actual %s", tt.envVars[i].Name, container.Env[i].Name)
					}
					if container.Env[i].Value != tt.envVars[i].Value {
						t.Errorf("TestGetContainer() error: expected env value %s, actual %s", tt.envVars[i].Value, container.Env[i].Value)
					}
				}
			}

			if len(container.Ports) != len(tt.ports) {
				t.Errorf("TestGetContainer() error: expected container port length %d, actual %d", len(tt.ports), len(container.Ports))
			} else {
				for i := range container.Ports {
					if container.Ports[i].Name != tt.ports[i].Name {
						t.Errorf("TestGetContainer() error: expected port name %s, actual %s", tt.ports[i].Name, container.Ports[i].Name)
					}
					if container.Ports[i].ContainerPort != tt.ports[i].ContainerPort {
						t.Errorf("TestGetContainer() error: expected port number %v, actual %v", tt.ports[i].ContainerPort, container.Ports[i].ContainerPort)
					}
				}
			}

		})
	}
}

func TestGetPodTemplateSpec(t *testing.T) {

	container := []corev1.Container{
		{
			Name:            "container1",
			Image:           "image1",
			ImagePullPolicy: corev1.PullAlways,

			Command: []string{"tail"},
			Args:    []string{"-f", "/dev/null"},
			Env:     []corev1.EnvVar{},
		},
	}

	volume := []corev1.Volume{
		{
			Name: "vol1",
		},
	}

	tests := []struct {
		podName        string
		namespace      string
		serviceAccount string
		labels         map[string]string
	}{
		{
			podName:        "podSpecTest",
			namespace:      "default",
			serviceAccount: "default",
			labels: map[string]string{
				"app":       "app",
				"component": "frontend",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.podName, func(t *testing.T) {

			objectMeta := GetObjectMeta(tt.podName, tt.namespace, tt.labels, nil)
			podTemplateSpecParams := podTemplateSpecParams{
				ObjectMeta:     objectMeta,
				Containers:     container,
				Volumes:        volume,
				InitContainers: container,
			}
			podTemplateSpec := getPodTemplateSpec(podTemplateSpecParams)

			if podTemplateSpec.Name != tt.podName {
				t.Errorf("TestGetPodTemplateSpec() error: expected podName %s, actual %s", tt.podName, podTemplateSpec.Name)
			}
			if podTemplateSpec.Namespace != tt.namespace {
				t.Errorf("TestGetPodTemplateSpec() error: expected namespace %s, actual %s", tt.namespace, podTemplateSpec.Namespace)
			}
			if !hasVolumeWithName("vol1", podTemplateSpec.Spec.Volumes) {
				t.Errorf("TestGetPodTemplateSpec() error: volume with name: %s not found", "vol1")
			}
			if !reflect.DeepEqual(podTemplateSpec.Labels, tt.labels) {
				t.Errorf("TestGetPodTemplateSpec() error: expected labels %+v, actual %+v", tt.labels, podTemplateSpec.Labels)
			}
			if !reflect.DeepEqual(podTemplateSpec.Spec.Containers, container) {
				t.Errorf("TestGetPodTemplateSpec() error: expected container %+v, actual %+v", container, podTemplateSpec.Spec.Containers)
			}
			if !reflect.DeepEqual(podTemplateSpec.Spec.InitContainers, container) {
				t.Errorf("TestGetPodTemplateSpec() error: expected InitContainers %+v, actual %+v", container, podTemplateSpec.Spec.InitContainers)
			}
		})
	}
}

func TestGetServiceSpec(t *testing.T) {

	endpointNames := []string{"port-8080-url", "port-9090-url", "a-very-long-port-name-before-endpoint-length-limit-8080"}

	tests := []struct {
		name                string
		containerComponents []v1.Component
		filteredComponents  []v1.Component
		labels              map[string]string
		filterOptions       common.DevfileOptions
		wantPorts           []corev1.ServicePort
	}{
		{
			name: "multiple endpoints have different ports",
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Endpoints: []v1.Endpoint{
								{
									Name:       endpointNames[0],
									TargetPort: 8080,
								},
								{
									Name:       endpointNames[1],
									TargetPort: 9090,
								},
							},
						},
					},
				},
			},
			labels: map[string]string{
				"component": "testcomponent",
			},
			wantPorts: []corev1.ServicePort{
				{
					Name:       endpointNames[0],
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
				{
					Name:       endpointNames[1],
					Port:       9090,
					TargetPort: intstr.FromInt(9090),
				},
			},
		},
		{
			name: "long port name before endpoint length limit to <=15",
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Endpoints: []v1.Endpoint{
								{
									Name:       endpointNames[2],
									TargetPort: 8080,
								},
							},
						},
					},
				},
			},
			labels: map[string]string{
				"component": "testcomponent",
			},
			wantPorts: []corev1.ServicePort{
				{
					Name:       "port-8080",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
		{
			name: "filter components",
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Endpoints: []v1.Endpoint{
								{
									Name:       endpointNames[0],
									TargetPort: 8080,
								},
							},
						},
					},
				},
				{
					Name: "testcontainer2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Endpoints: []v1.Endpoint{
								{
									Name:       endpointNames[1],
									TargetPort: 9090,
								},
							},
						},
					},
				},
			},
			labels: map[string]string{
				"component": "testcomponent",
			},
			wantPorts: []corev1.ServicePort{
				{
					Name:       endpointNames[1],
					Port:       9090,
					TargetPort: intstr.FromInt(9090),
				},
			},
			filteredComponents: []v1.Component{
				{
					Name: "testcontainer2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Endpoints: []v1.Endpoint{
								{
									Name:       endpointNames[1],
									TargetPort: 9090,
								},
							},
						},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString": "firstStringValue",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDevfileData := data.NewMockDevfileData(ctrl)

			tt.filterOptions.ComponentOptions = common.ComponentOptions{
				ComponentType: v1.ContainerComponentType,
			}
			mockGetComponents := mockDevfileData.EXPECT().GetComponents(tt.filterOptions)

			// set up the mock data
			if len(tt.filterOptions.Filter) == 0 {
				mockGetComponents.Return(tt.containerComponents, nil).AnyTimes()
			} else {
				mockGetComponents.Return(tt.filteredComponents, nil).AnyTimes()
			}
			mockDevfileData.EXPECT().GetProjects(common.DevfileOptions{}).Return(nil, nil).AnyTimes()
			mockDevfileData.EXPECT().GetEvents().Return(v1.Events{}).AnyTimes()
			devObj := parser.DevfileObj{
				Data: mockDevfileData,
			}

			serviceSpec, err := getServiceSpec(devObj, tt.labels, tt.filterOptions)

			// Unexpected error
			if err != nil {
				t.Errorf("TestGetServiceSpec() unexpected error: %v", err)
			} else {
				if !reflect.DeepEqual(serviceSpec.Selector, tt.labels) {
					t.Errorf("TestGetServiceSpec() error: expected service selector is %v, actual %v", tt.labels, serviceSpec.Selector)
				}
				if len(serviceSpec.Ports) != len(tt.wantPorts) {
					t.Errorf("TestGetServiceSpec() error: expected service ports length is %v, actual %v", len(tt.wantPorts), len(serviceSpec.Ports))
				} else {
					for i := range serviceSpec.Ports {
						if serviceSpec.Ports[i].Name != tt.wantPorts[i].Name {
							t.Errorf("TestGetServiceSpec() error: expected name %s, actual name %s", tt.wantPorts[i].Name, serviceSpec.Ports[i].Name)
						}
						if serviceSpec.Ports[i].Port != tt.wantPorts[i].Port {
							t.Errorf("TestGetServiceSpec() error: expected port number is %v, actual %v", tt.wantPorts[i].Port, serviceSpec.Ports[i].Port)
						}
					}
				}
			}
		})
	}
}

func TestGetPortExposure(t *testing.T) {
	urlName := "testurl"
	urlName2 := "testurl2"
	tests := []struct {
		name                string
		containerComponents []v1.Component
		filteredComponents  []v1.Component
		filterOptions       common.DevfileOptions
		wantMap             map[int]v1.EndpointExposure
	}{
		{
			name: "devfile has single container with single endpoint",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "devfile no endpoints",
			wantMap: map[int]v1.EndpointExposure{},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
						},
					},
				},
			},
		},
		{
			name: "devfile has multiple endpoints with same port, 1 public and 1 internal, should assign public",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.InternalEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "devfile has multiple endpoints with same port, 1 public and 1 none, should assign public",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "devfile has multiple endpoints with same port, 1 internal and 1 none, should assign internal",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.InternalEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.InternalEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "devfile has multiple endpoints with different port",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
				9090: v1.InternalEndpointExposure,
				3000: v1.NoneEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
								},
								{
									Name:       urlName,
									TargetPort: 3000,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
				{
					Name: "testcontainer2",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName2,
									TargetPort: 9090,
									Secure:     &isTrue,
									Path:       "/testpath",
									Exposure:   v1.InternalEndpointExposure,
									Protocol:   v1.HTTPSEndpointProtocol,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Filter components",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
				3000: v1.NoneEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
								},
								{
									Name:       urlName,
									TargetPort: 3000,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
				{
					Name: "testcontainer2",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName2,
									TargetPort: 9090,
									Secure:     &isTrue,
									Path:       "/testpath",
									Exposure:   v1.InternalEndpointExposure,
									Protocol:   v1.HTTPSEndpointProtocol,
								},
							},
						},
					},
				},
			},
			filteredComponents: []v1.Component{
				{
					Name: "testcontainer1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
								},
								{
									Name:       urlName,
									TargetPort: 3000,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString": "firstStringValue",
				},
			},
		},
		{
			name:    "Wrong filter components",
			wantMap: map[int]v1.EndpointExposure{},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
								},
								{
									Name:       urlName,
									TargetPort: 3000,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
			filteredComponents: nil,
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstStringWrong": "firstStringValue",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDevfileData := data.NewMockDevfileData(ctrl)

			tt.filterOptions.ComponentOptions = common.ComponentOptions{
				ComponentType: v1.ContainerComponentType,
			}
			mockGetComponents := mockDevfileData.EXPECT().GetComponents(tt.filterOptions)

			// set up the mock data
			if len(tt.filterOptions.Filter) == 0 {
				mockGetComponents.Return(tt.containerComponents, nil).AnyTimes()
			} else {
				mockGetComponents.Return(tt.filteredComponents, nil).AnyTimes()
			}
			devObj := parser.DevfileObj{
				Data: mockDevfileData,
			}

			mapCreated, err := getPortExposure(devObj, tt.filterOptions)
			// Checks for unexpected error cases
			if err != nil {
				t.Errorf("TestGetPortExposure() unexpected error: %v", err)
			} else if !reflect.DeepEqual(mapCreated, tt.wantMap) {
				t.Errorf("TestGetPortExposure() error: expected: %v, got %v", tt.wantMap, mapCreated)
			}

		})
	}

}

func TestGetIngressSpec(t *testing.T) {

	tests := []struct {
		name      string
		parameter IngressSpecParams
	}{
		{
			name: "1",
			parameter: IngressSpecParams{
				ServiceName:   "service1",
				IngressDomain: "test.1.2.3.4.nip.io",
				PortNumber: intstr.IntOrString{
					IntVal: 8080,
				},
				TLSSecretName: "testTLSSecret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ingressSpec := getIngressSpec(tt.parameter)

			if ingressSpec.Rules[0].Host != tt.parameter.IngressDomain {
				t.Errorf("TestGetIngressSpec() error: expected ingressDomain %s, actual %s", tt.parameter.IngressDomain, ingressSpec.Rules[0].Host)
			}

			if ingressSpec.Rules[0].HTTP.Paths[0].Backend.ServicePort != tt.parameter.PortNumber {
				t.Errorf("TestGetIngressSpec() error: expected portNumber %v, actual %v", tt.parameter.PortNumber, ingressSpec.Rules[0].HTTP.Paths[0].Backend.ServicePort)
			}

			if ingressSpec.Rules[0].HTTP.Paths[0].Backend.ServiceName != tt.parameter.ServiceName {
				t.Errorf("TestGetIngressSpec() error: expected serviceName %s, actual %s", tt.parameter.ServiceName, ingressSpec.Rules[0].HTTP.Paths[0].Backend.ServiceName)
			}

			if ingressSpec.TLS[0].SecretName != tt.parameter.TLSSecretName {
				t.Errorf("TestGetIngressSpec() error: expected TLSSecretName %s, actual %s", tt.parameter.TLSSecretName, ingressSpec.TLS[0].SecretName)
			}

		})
	}
}

func TestGetNetworkingV1IngressSpec(t *testing.T) {

	tests := []struct {
		name      string
		parameter IngressSpecParams
	}{
		{
			name: "1",
			parameter: IngressSpecParams{
				ServiceName:   "service1",
				IngressDomain: "test.1.2.3.4.nip.io",
				PortNumber: intstr.IntOrString{
					IntVal: 8080,
				},
				TLSSecretName: "testTLSSecret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ingressSpec := getNetworkingV1IngressSpec(tt.parameter)

			if ingressSpec.Rules[0].Host != tt.parameter.IngressDomain {
				t.Errorf("TestGetNetworkingV1IngressSpec() error: expected IngressDomain %s, actual %s", tt.parameter.IngressDomain, ingressSpec.Rules[0].Host)
			}

			if ingressSpec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number != tt.parameter.PortNumber.IntVal {
				t.Errorf("TestGetNetworkingV1IngressSpec() error: expected PortNumber %v, actual %v", tt.parameter.PortNumber, ingressSpec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
			}

			if ingressSpec.Rules[0].HTTP.Paths[0].Backend.Service.Name != tt.parameter.ServiceName {
				t.Errorf("TestGetNetworkingV1IngressSpec() error: expected ServiceName %s, actual %s", tt.parameter.ServiceName, ingressSpec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
			}

			if ingressSpec.TLS[0].SecretName != tt.parameter.TLSSecretName {
				t.Errorf("TestGetNetworkingV1IngressSpec() error: expected TLSSecretName %s, actual %s", tt.parameter.TLSSecretName, ingressSpec.TLS[0].SecretName)
			}

		})
	}
}

func TestGetRouteSpec(t *testing.T) {

	tests := []struct {
		name      string
		parameter RouteSpecParams
	}{
		{
			name: "insecure route",
			parameter: RouteSpecParams{
				ServiceName: "service1",
				PortNumber: intstr.IntOrString{
					IntVal: 8080,
				},
				Secure: false,
				Path:   "/test",
			},
		},
		{
			name: "secure route",
			parameter: RouteSpecParams{
				ServiceName: "service1",
				PortNumber: intstr.IntOrString{
					IntVal: 8080,
				},
				Secure: true,
				Path:   "/test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			routeSpec := getRouteSpec(tt.parameter)

			if routeSpec.Port.TargetPort != tt.parameter.PortNumber {
				t.Errorf("TestGetRouteSpec() error: expected PortNumber %v, actual %v", tt.parameter.PortNumber, routeSpec.Port.TargetPort)
			}

			if routeSpec.To.Name != tt.parameter.ServiceName {
				t.Errorf("TestGetRouteSpec() error: expected ServiceName %s, actual %s", tt.parameter.ServiceName, routeSpec.To.Name)
			}

			if routeSpec.Path != tt.parameter.Path {
				t.Errorf("TestGetRouteSpec() error: expected Path %s, actual %s", tt.parameter.Path, routeSpec.Path)
			}

			if (routeSpec.TLS != nil) != tt.parameter.Secure {
				t.Errorf("TestGetRouteSpec() error: the route TLS does not match secure level %v", tt.parameter.Secure)
			}

		})
	}
}

func TestGetPVCSpec(t *testing.T) {

	tests := []struct {
		name    string
		size    string
		wantErr bool
	}{
		{
			name:    "Valid resource size",
			size:    "1Gi",
			wantErr: false,
		},
		{
			name:    "Resource size missing",
			size:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			quantity, err := resource.ParseQuantity(tt.size)
			// Checks for unexpected error cases
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestGetPVCSpec() error: resource.ParseQuantity unexpected error %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				pvcSpec := getPVCSpec(quantity)
				if pvcSpec.AccessModes[0] != corev1.ReadWriteOnce {
					t.Errorf("TestGetPVCSpec() error: AccessMode Error: expected %s, actual %s", corev1.ReadWriteMany, pvcSpec.AccessModes[0])
				}

				pvcSpecQuantity := pvcSpec.Resources.Requests["storage"]
				if pvcSpecQuantity.String() != quantity.String() {
					t.Errorf("TestGetPVCSpec() error: pvcSpec.Resources.Requests Error: expected %v, actual %v", pvcSpecQuantity.String(), quantity.String())
				}
			}
		})
	}
}

func hasVolumeWithName(name string, volMounts []corev1.Volume) bool {
	for _, vm := range volMounts {
		if vm.Name == name {
			return true
		}
	}
	return false
}

func TestGetBuildConfigSpec(t *testing.T) {

	image := "image"
	namespace := "namespace"

	tests := []struct {
		name          string
		GitURL        string
		GitRef        string
		ContextDir    string
		buildStrategy buildv1.BuildStrategy
	}{
		{
			name:          "Get a Source Strategy Build Config",
			GitURL:        "url",
			GitRef:        "ref",
			buildStrategy: GetSourceBuildStrategy(image, namespace),
		},
		{
			name:          "Get a Docker Strategy Build Config",
			GitURL:        "url",
			GitRef:        "ref",
			ContextDir:    "./",
			buildStrategy: GetDockerBuildStrategy("dockerfilePath", []corev1.EnvVar{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commonObjectMeta := GetObjectMeta(image, namespace, nil, nil)
			params := BuildConfigSpecParams{
				ImageStreamTagName: commonObjectMeta.Name,
				BuildStrategy:      tt.buildStrategy,
				GitURL:             tt.GitURL,
				GitRef:             tt.GitRef,
				ContextDir:         tt.ContextDir,
			}
			buildConfigSpec := getBuildConfigSpec(params)

			if !strings.Contains(buildConfigSpec.CommonSpec.Output.To.Name, image) {
				t.Error("TestGetBuildConfigSpec() error: build config output name does not match")
			}

			if buildConfigSpec.Source.Git.Ref != tt.GitRef || buildConfigSpec.Source.Git.URI != tt.GitURL {
				t.Error("TestGetBuildConfigSpec() error: build config git source does not match")
			}

			if buildConfigSpec.CommonSpec.Source.ContextDir != tt.ContextDir {
				t.Error("TestGetBuildConfigSpec() error: context dir does not match")
			}
		})
	}

}

func TestGetPVC(t *testing.T) {

	tests := []struct {
		name       string
		pvc        string
		volumeName string
	}{
		{
			name:       "Get PVC vol for given pvc name and volume name",
			pvc:        "mypvc",
			volumeName: "myvolume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume := getPVC(tt.volumeName, tt.pvc)

			if volume.Name != tt.volumeName {
				t.Errorf("TestGetPVC() error: volume name does not match; expected %s got %s", tt.volumeName, volume.Name)
			}

			if volume.PersistentVolumeClaim.ClaimName != tt.pvc {
				t.Errorf("TestGetPVC() error: pvc name does not match; expected %s got %s", tt.pvc, volume.PersistentVolumeClaim.ClaimName)
			}
		})
	}
}

func TestAddVolumeMountToContainers(t *testing.T) {

	tests := []struct {
		name                   string
		volumeName             string
		containerMountPathsMap map[string][]string
		container              corev1.Container
	}{
		{
			name:       "Successfully mount volume to container",
			volumeName: "myvolume",
			containerMountPathsMap: map[string][]string{
				"container1": {"/tmp/path1", "/tmp/path2"},
			},
			container: corev1.Container{
				Name:            "container1",
				Image:           "image1",
				ImagePullPolicy: corev1.PullAlways,

				Command: []string{"tail"},
				Args:    []string{"-f", "/dev/null"},
				Env:     []corev1.EnvVar{},
			},
		},
		{
			name:       "No Container present to mount volume",
			volumeName: "myvolume",
			containerMountPathsMap: map[string][]string{
				"container1": {"/tmp/path1", "/tmp/path2"},
			},
			container: corev1.Container{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containers := []corev1.Container{tt.container}
			addVolumeMountToContainers(containers, tt.volumeName, tt.containerMountPathsMap)

			mountPathCount := 0
			for _, container := range containers {
				if container.Name == tt.container.Name {
					for _, volumeMount := range container.VolumeMounts {
						if volumeMount.Name == tt.volumeName {
							for _, mountPath := range tt.containerMountPathsMap[tt.container.Name] {
								if volumeMount.MountPath == mountPath {
									mountPathCount++
								}
							}
						}
					}
				}
			}

			if mountPathCount != len(tt.containerMountPathsMap[tt.container.Name]) {
				t.Errorf("TestAddVolumeMountToContainers() error: Volume Mounts for %s have not been properly mounted to the container", tt.volumeName)
			}
		})
	}
}

func TestGetContainerAnnotations(t *testing.T) {
	trueBool := true

	tests := []struct {
		name                string
		containerComponents []v1.Component
		expected            v1.Annotation
	}{
		{
			name: "no dedicated pod",
			containerComponents: []v1.Component{
				testingutil.GenerateDummyContainerComponent("container1", nil, nil, nil, v1.Annotation{
					Service: map[string]string{
						"key1": "value1",
					},
					Deployment: map[string]string{
						"key1": "value1",
					},
				}, nil),
				testingutil.GenerateDummyContainerComponent("container2", nil, nil, nil, v1.Annotation{
					Service: map[string]string{
						"key2": "value2",
					},
					Deployment: map[string]string{
						"key2": "value2",
					},
				}, nil),
			},
			expected: v1.Annotation{
				Service: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Deployment: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		{
			name: "has dedicated pod",
			containerComponents: []v1.Component{
				testingutil.GenerateDummyContainerComponent("container1", nil, nil, nil, v1.Annotation{
					Service: map[string]string{
						"key1": "value1",
					},
					Deployment: map[string]string{
						"key1": "value1",
					},
				}, nil),
				testingutil.GenerateDummyContainerComponent("container2", nil, nil, nil, v1.Annotation{
					Service: map[string]string{
						"key2": "value2",
					},
					Deployment: map[string]string{
						"key2": "value2",
					},
				}, &trueBool),
			},
			expected: v1.Annotation{
				Service: map[string]string{
					"key1": "value1",
				},
				Deployment: map[string]string{
					"key1": "value1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDevfileData := data.NewMockDevfileData(ctrl)

			options := common.DevfileOptions{
				ComponentOptions: common.ComponentOptions{
					ComponentType: v1.ContainerComponentType,
				},
			}
			// set up the mock data
			mockGetComponents := mockDevfileData.EXPECT().GetComponents(options)
			mockGetComponents.Return(tt.containerComponents, nil).AnyTimes()

			devObj := parser.DevfileObj{
				Data: mockDevfileData,
			}
			annotations, err := getContainerAnnotations(devObj, common.DevfileOptions{})
			// Checks for unexpected error cases
			if err != nil {
				t.Errorf("TestGetContainerAnnotations(): unexpected error %v", err)
			}
			assert.Equal(t, tt.expected, annotations, "TestGetContainerAnnotations(): The two values should be the same.")

		})
	}
}

func TestMergeMaps(t *testing.T) {

	tests := []struct {
		name     string
		dest     map[string]string
		src      map[string]string
		expected map[string]string
	}{
		{
			name: "dest is nil",
			dest: nil,
			src: map[string]string{
				"key3": "value3",
			},
			expected: map[string]string{
				"key3": "value3",
			},
		},
		{
			name: "src is nil",
			dest: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			src: nil,
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "no nil maps",
			dest: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			src: map[string]string{
				"key3": "value3",
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeMaps(tt.dest, tt.src)
			assert.Equal(t, tt.expected, result, "TestmergeMaps(): The two values should be the same.")

		})
	}
}
