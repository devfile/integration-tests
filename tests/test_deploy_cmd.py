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
            result = subprocess.run(["odo", "deploy", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "build -t quay.io/unknown-account/myimage -f "
                            + os.path.abspath(os.path.join(self.CONTEXT, "Dockerfile"))
                            + " "
                            + os.path.abspath(self.CONTEXT))
            assert contains(result.stdout, "push quay.io/unknown-account/myimage")

            result = subprocess.run(["kubectl", "get", "deployment", "my-component", "-n", "intg-test-project",
                                     "-o", "jsonpath='{.spec.template.spec.containers[0].image}'"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "quay.io/unknown-account/myimage")
