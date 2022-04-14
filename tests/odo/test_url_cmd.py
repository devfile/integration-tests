import tempfile
import jmespath

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestUrlCmd:

    CONTAINER_NAME = "test-container"
    CONTEXT = "test-context"
    HOST = "test.host.com"
    PORT_1 = "3000"
    PORT_2 = "5000"
    ENDPOINT_1 = "url-1"
    ENDPOINT_2 = "url-2"

    ENDPOINT = "3000-tcp"

    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])


    def test_url_duplicate_name_port(self):
        print("Test case : should not allow to create endpoint with duplicate name or port")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            # should not allow creating an endpoint with duplicate name
            result = subprocess.run(["odo", "url", "create", self.ENDPOINT, "--port", self.PORT_1,
                                     "--host", self.HOST, "--secure", "--ingress"],
                                    capture_output=True, text=True, check=False)
            assert contains(result.stderr,
                            "url {} already exist in devfile endpoint entry under container runtime".format(self.ENDPOINT))

            # should not allow to create URL with duplicate port
            result = subprocess.run(["odo", "url", "create", self.ENDPOINT_1, "--port", self.PORT_1,
                                     "--host", self.HOST, "--secure", "--ingress"],
                                    capture_output=True, text=True, check=False)

            # Todo: potential bug - it's not blocked by the odo used in the test. Need to verify if it's fixed in more recent release
            # assert contains(result.stdout, "port 3000 already exists in devfile endpoint entry")

    def test_url_invalid_container(self):
        print("Test case : should not allow creating under an invalid container")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            # should not allow creating an endpoint with duplicate name
            result = subprocess.run(["odo", "url", "create", "--port", self.PORT_1, "--host", self.HOST,
                                     "--container", self.CONTAINER_NAME, "--ingress"],
                                    capture_output=True, text=True, check=False)
            assert contains(result.stderr,
                            "the container specified: {} does not exist in devfile".format(self.CONTAINER_NAME))

    def test_url_create_and_delete(self):

        print('Test case : should successfully create endpoint')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            subprocess.run(["odo", "url", "create", self.ENDPOINT_1, "--port", self.PORT_1,
                            "--host", self.HOST, "--secure", "--ingress"])

            list_components_before_push = [
                self.ENDPOINT_1,
                "Not Pushed",
                "true",
                "ingress"
            ]

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=True)

            assert match_all(result.stdout, list_components_before_push)

            subprocess.run(["odo", "url", "delete", self.ENDPOINT_1, "-f"],
                                    capture_output=True, text=True, check=True)

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=False)

            assert contains(result.stderr, "no URLs found for component nodejs")


    def test_url_create_multiple_endpoints(self):

        print('Test case : should successfully create multiple endpoints')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            # create 1st url
            subprocess.run(["odo", "url", "create", "--port", self.PORT_1,
                            "--host", self.HOST, "--secure", "--ingress"])

            list_components_after_url1 = [
                self.ENDPOINT,
                self.PORT_1,
                "Not Pushed",
                "true",
                "ingress"
            ]

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=True)

            assert match_all(result.stdout, list_components_after_url1)

            # creat 2nd url
            subprocess.run(["odo", "url", "create", self.ENDPOINT_2, "--port", self.PORT_2,
                            "--host", self.HOST, "--secure", "--ingress"])

            list_components_after_url2 = [
                self.ENDPOINT_2,
                self.PORT_2,
                "Not Pushed",
                "true",
                "ingress"
            ]

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=True)

            assert match_all(result.stdout, list_components_after_url2)

            # odo url delete 1st one
            subprocess.run(["odo", "url", "delete", self.ENDPOINT, "-f"],
                                    capture_output=True, text=True, check=True)

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, self.ENDPOINT) == False

            # odo url delete 2nd one
            subprocess.run(["odo", "url", "delete", self.ENDPOINT_2, "-f"],
                           capture_output=True, text=True, check=True)

            result = subprocess.run(["odo", "url", "list"],
                                    capture_output=True, text=True, check=False)

            assert contains(result.stderr, "no URLs found for component nodejs")
