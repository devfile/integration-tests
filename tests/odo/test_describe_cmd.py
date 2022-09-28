import tempfile
import jmespath

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestDescribeCmd:

    APP_NAME = "test-app"
    COMP_NAME = "cmp-django"
    COMP_TYPE = "Django"
    CONTEXT = "test-context"
    PORT_1 = "3000"
    PORT_2 = "4000"
    SIZE = "1Gi"
    STORAGE_PATH = "/data1"
    STORAGE_1 = "storage-1"
    ENDPOINT_1 = "url-1"
    ENDPOINT_2 = "url-2"
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

    def test_describe_components(self):
        print('Test case : should describe the component correctly')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            subprocess.run(["odo", "create", "python-django", self.COMP_NAME, "--project", self.tmp_project_name,
                            "--context", self.CONTEXT, "--app", self.APP_NAME])

            subprocess.run(["odo", "url", "create", self.ENDPOINT_1, "--port", self.PORT_1, "--host", self.HOST, "--context", self.CONTEXT])
            subprocess.run(["odo", "url", "create", self.ENDPOINT_2, "--port", self.PORT_2, "--host", self.HOST, "--context", self.CONTEXT])
            subprocess.run(["odo", "storage", "create", self.STORAGE_1, "--size", self.SIZE, "--path", self.STORAGE_PATH,
                            "--context", self.CONTEXT])

            result = subprocess.run(["odo", "describe", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)

            list_components = [
                self.COMP_NAME,
                self.COMP_TYPE,
                self.ENDPOINT_1,
                self.ENDPOINT_2,
                self.STORAGE_1,
            ]

            assert match_all(result.stdout, list_components)

            result = subprocess.run(["odo", "describe", self.COMP_NAME, "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)
            assert(contains(result.stdout, self.COMP_NAME))

            result = subprocess.run(["odo", "describe", "-o", "json", "--context", self.CONTEXT],
                                    capture_output=True, text=True, check=True)
            assert validate_json_format(result.stdout)

            print("odo describe -o json --context:", result.stdout)

            dict = json.loads(result.stdout)
            path = jmespath.search('metadata.name', dict)
            assert contains(path, self.COMP_NAME)

            path = jmespath.search('metadata.namespace', dict)
            assert contains(path, self.tmp_project_name)

            path = jmespath.search('spec.app', dict)
            assert contains(path, self.APP_NAME)

            path = jmespath.search('spec.type', dict)
            assert contains(path, self.COMP_TYPE)

            path = jmespath.search('spec.urls.items[0].metadata.name', dict)
            assert contains(path, self.ENDPOINT_1)

            path = jmespath.search('spec.urls.items[0].spec.port', dict)
            assert contains(str(path), self.PORT_1)

            path = jmespath.search('spec.urls.items[1].metadata.name', dict)
            assert contains(path, self.ENDPOINT_2)

            path = jmespath.search('spec.urls.items[1].spec.port', dict)
            assert contains(str(path), self.PORT_2)

            path = jmespath.search('spec.storages.items[0].metadata.name', dict)
            assert contains(path, self.STORAGE_1)

            path = jmespath.search('spec.storages.items[0].spec.size', dict)
            assert contains(path, self.SIZE)

            path = jmespath.search('spec.storages.items[0].spec.path', dict)
            assert contains(path, self.STORAGE_PATH)

            path = jmespath.search('status.state', dict)
            assert contains(path, "Not Pushed")

            # Todo:
            # result = subprocess.run(["odo", "push", "--context", self.CONTEXT],
            #                         capture_output=True, text=True, check=True)
            #
            # dict = json.loads(result.stdout)
            # path = jmespath.search('status.state', dict)
            # assert contains(path, "Pushed")

    def test_describe_missing_metadata(self):
        print("Test case : when 'projectType' is missing, it should show the language for 'Type' in odo describe")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/springboot/devfile-with-missing-projectType-metadata.yaml')

            copy_and_create(source_devfile_path, "springboot/project", tmp_workspace, self.CONTEXT)

            result = subprocess.run(["odo", "describe", "--context", os.path.join(tmp_workspace, self.CONTEXT), "-o", "json"],
                                    capture_output=True, text=True, check=True)

            assert validate_json_format(result.stdout)

            dict = json.loads(result.stdout)
            path = jmespath.search('spec.type', dict)
            assert contains(path, "java")

    def test_describe_missing_metadata_and_language(self):
        print("Test case : when both 'projectType' and 'language' are missing, it should show 'Not available' for 'Type' in odo describe")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/springboot/devfile-with-missing-projectType-and-language-metadata.yaml')

            copy_and_create(source_devfile_path, "springboot/project", tmp_workspace, self.CONTEXT)

            result = subprocess.run(
                ["odo", "describe", "--context", os.path.join(tmp_workspace, self.CONTEXT), "-o", "json"],
                capture_output=True, text=True, check=True)

            assert validate_json_format(result.stdout)

            dict = json.loads(result.stdout)
            path = jmespath.search('spec.type', dict)
            assert contains(path, "Not available")
