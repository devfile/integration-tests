import tempfile

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestBuildImagesCmd:

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

    def test_build_push_images(self):

        print("Test case : should run odo build-images and odo push by using a devfile.yaml containing an Image component")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-outerloop.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"
            result = subprocess.run(["odo", "build-images", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "build -t quay.io/unknown-account/myimage -f "
                            + os.path.abspath(os.path.join(self.CONTEXT, "Dockerfile"))
                            + " "
                            + os.path.abspath(self.CONTEXT))

            result = subprocess.run(["odo", "build-images", "--context", self.CONTEXT, "--push"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "push quay.io/unknown-account/myimage")

    def test_build_images_with_dockerfile_args(self):
        print("Test case : should run odo build-images by using a devfile.yaml containing an Image component with Dockerfile args")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-outerloop-args.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"
            result = subprocess.run(["odo", "build-images", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "--unknown-flag value")
            