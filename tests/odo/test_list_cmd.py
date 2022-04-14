import tempfile
import jmespath

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestListCmd:

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


    def test_list_component_in_application(self):

        print("Test case : should successfully list component when a component is created in 'app' application")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            list_components = [
                "app",
                "nodejs",
                "Not Pushed"
            ]

            result = subprocess.run(["odo", "list"],
                                    capture_output=True, text=True, check=True)

            assert match_all(result.stdout, list_components)

            result_json = subprocess.run(["odo", "list", "-o", "json"],
                                    capture_output=True, text=True, check=True)

            assert validate_json_format(result_json.stdout)

            dict = json.loads(result_json.stdout)

            path = jmespath.search('kind', dict)
            assert contains(path, "List")
            path = jmespath.search('devfileComponents[0].kind', dict)
            assert contains(path, "Component")
            path = jmespath.search('devfileComponents[0].metadata.name', dict)
            assert contains(path, "nodejs")
            path = jmespath.search('devfileComponents[0].status.state', dict)
            assert contains(path, "Not Pushed")

    def test_list_component_missing_metadata_projecttype(self):

        print("Test case : should show the language for 'Type' in odo list")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/springboot/devfile-with-missing-projectType-metadata.yaml')

            copy_and_create(source_devfile_path, "springboot/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            list_components = [
                "app",
                "java",
                "myspringbootproject",
                "Not Pushed"
            ]

            result = subprocess.run(["odo", "list"],
                                    capture_output=True, text=True, check=True)
            assert match_all(result.stdout, list_components)

            result_json = subprocess.run(["odo", "list", "-o", "json"],
                                    capture_output=True, text=True, check=True)
            assert validate_json_format(result_json.stdout)

            dict = json.loads(result_json.stdout)
            path = jmespath.search('devfileComponents[0].spec.type', dict)
            assert contains(path, "java")


    def test_list_component_missing_metadata_projecttype_language(self):

        print("Test case : should show 'Not available' for 'Type' in odo list")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/springboot/devfile-with-missing-projectType-and-language-metadata.yaml')

            copy_and_create(source_devfile_path, "springboot/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "list"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Not available")
