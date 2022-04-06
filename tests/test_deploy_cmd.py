import sys
import tempfile

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestDeployCmd:

    CONTEXT = "test-context"

    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])

    def test_deploy(self):

        print("Test case : should run odo deploy by using a devfile.yaml containing a deploy command")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-deploy.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"

            # building and pushing image to registry
            result = subprocess.run(["odo", "deploy", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "build -t quay.io/unknown-account/myimage -f "
                            + os.path.abspath(os.path.join(self.CONTEXT, "Dockerfile"))
                            + " "
                            + os.path.abspath(self.CONTEXT))
            assert contains(result.stdout, "push quay.io/unknown-account/myimage")

            # deploying a deployment with the built image
            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(["minikube", "kubectl", "--", "get", "deployment", "my-component", "-n", "intg-test-project",
                                     "-o", "jsonpath='{.spec.template.spec.containers[0].image}'"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "quay.io/unknown-account/myimage")

    def test_deploy_image_component(self):

        print("Test case : should run odo deploy by using a devfile.yaml containing an Image component with a build context")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-outerloop-project_source-in-docker-build-context.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"

            # building and pushing image to registry
            result = subprocess.run(["odo", "deploy", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "build -t localhost:5000/devfile-nodejs-deploy:0.1.0 -f "
                            + os.path.abspath(os.path.join(self.CONTEXT, "Dockerfile")))

            assert contains(result.stdout, "push localhost:5000/devfile-nodejs-deploy:0.1.0")
            assert contains(result.stdout, "Deploying Kubernetes Deployment: devfile-nodejs-deploy")
            assert contains(result.stdout, "Deploying Kubernetes Service: devfile-nodejs-deploy")


            # deploying a deployment with the built image
            # MacOS: reuse the existing kubectl inside minikube
            # result = subprocess.run(["minikube", "kubectl", "--", "get", "deployment", "my-component", "-n", "intg-test-project",
            #                          "-o", "jsonpath='{.spec.template.spec.containers[0].image}'"],
            #                         capture_output=True, text=True, check=True)
            # assert contains(result.stdout, "quay.io/unknown-account/myimage")


    # Supporting multiple deploy commands are not supported yet
    # issue: https://github.com/redhat-developer/odo/issues/5454
    # @pytest.mark.skip(reason=None)
    # def test_two_deploy_commands(self):
    #
    #     print("Test case : should run odo deploy by using a devfile.yaml containing a deploy command")
    #
    #     with tempfile.TemporaryDirectory() as tmp_workspace:
    #         os.chdir(tmp_workspace)
    #
    #         # example devfile path
    #         source_devfile_path = os.path.join(os.path.dirname(__file__),
    #                                            'examples/source/devfiles/nodejs/devfile-with-two-deploy-commands.yaml')
    #
    #         copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)
    #
    #         os.environ['PODMAN_CMD'] = "echo"
    #
    #         # building and pushing image to registry
    #         result = subprocess.run(["odo", "deploy", "--context", self.CONTEXT],
    #                                 capture_output=True, text=True, check=True)
    #
    #         assert contains(result.stdout, "build -t quay.io/unknown-account/myimage -f "
    #                         + os.path.abspath(os.path.join(self.CONTEXT, "Dockerfile"))
    #                         + " "
    #                         + os.path.abspath(self.CONTEXT))
    #         assert contains(result.stdout, "push quay.io/unknown-account/myimage")
    #
    #         # deploying a deployment with the built image
    #         # MacOS: reuse the existing kubectl inside minikube
    #         result = subprocess.run(["minikube", "kubectl", "--", "get", "deployment", "my-component", "-n", "intg-test-project",
    #                                  "-o", "jsonpath='{.spec.template.spec.containers[0].image}'"],
    #                                 capture_output=True, text=True, check=True)
    #         assert contains(result.stdout, "quay.io/unknown-account/myimage")
