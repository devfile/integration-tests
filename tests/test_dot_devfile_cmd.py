import tempfile

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestDotDevfile:

    CONTEXT = "test-context"
    PORT_1 = "9090"
    ENDPOINT_1 = "url-1"
    HOST = "test.host.com"

    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])


    def test_url_create_with_dot_devfile(self):

        print('Test case : should successfully create url by using .devfile.yaml')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            subprocess.run(["mv", "devfile.yaml", ".devfile.yaml"])
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

            # Todo:
            # list_components_after_push = [
            #     self.ENDPOINT_1,
            #     "Pushed",
            #     "true",
            #     "ingress"
            # ]
            #
            # subprocess.run(["odo", "push"])
            # result = subprocess.run(["odo", "url", "list"],
            #                         capture_output=True, text=True, check=True)
            #
            # assert match_all(result.stdout, list_components_after_push)

            # subprocess.run(["odo", "url", "delete", self.ENDPOINT_1, "-f"])

            # result = subprocess.run(["odo", "url", "list"],
            #                         capture_output=True, text=True, check=True)
            #
            # list_components_after_url_delete = [
            #     self.ENDPOINT_1,
            #     "Locally Deleted",
            #     "true",
            #     "ingress"
            # ]
            #
            # assert match_all(result.stdout, list_components_after_url_delete)
            # subprocess.run(["odo", "push"])

