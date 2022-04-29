import tempfile
import jmespath
import time

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestCreateCmd:
    tmp_project_name = None


    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()


    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])


    @pytest.mark.parametrize(
       "value, param_1, param_2, param_3", [pytest.param("2.1.0", "schemaVersion", -1,         -1, id = "check-schema-version"),
                                            pytest.param("java",  "metadata",   "language",    -1, id = "check-metadata-language"),
                                            pytest.param("dev",   "components",      0,    "name", id = "check-open-liberty-comp-name")])
    def test_create_component(self, value, param_1, param_2, param_3):

        print('Test case : should successfully create the devfile component and check contents')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            component_name = random_string()
            subprocess.run(["odo", "create", "java-openliberty", component_name])

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, param_1, param_2, param_3) == value, 'validate yaml contents'


    def test_create_component_with_valid_component_name(self):

        print('Test case : should successfully create the devfile component with valid component name')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            component_name = random_string()
            subprocess.run(["odo", "create", "java-openliberty", component_name])

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == component_name


    def test_create_component_with_project_flag(self):

        print('Test case : should successfully create the devfile component with --project flag')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            component_namespace = random_string()
            subprocess.run(["odo", "create", "java-openliberty", "--project", component_namespace])
            time.sleep(5)

            envfile_path = os.path.abspath(os.path.join(tmp_workspace, '.odo/env/env.yaml'))

            if os.path.isfile(envfile_path):
                assert query_yaml(envfile_path, "ComponentSettings", "Project", -1) == component_namespace
            else:
                raise ValueError("Failed: %s is not created yet." % file_path)


    def test_create_with_context_flag(self):
        print('Test case : odo create is executed with the --context flag')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            context_name = random_string()

            subprocess.run(["odo", "create", "nodejs", "--context", context_name])
            devfile_path = os.path.abspath(os.path.join(tmp_workspace, context_name, 'devfile.yaml'))
            envfile_path = os.path.abspath(os.path.join(tmp_workspace, context_name, '.odo/env/env.yaml'))

            assert os.path.exists(devfile_path)
            assert os.path.exists(envfile_path)

            result = subprocess.run(["odo", "env", "view", "--context", context_name, "-o", "json"],
                                    capture_output=True, text=True, check=True)

            assert validate_json_format(result.stdout)

            dict = json.loads(result.stdout)
            path = jmespath.search('spec.name', dict)
            assert contains(path, "nodejs-")


    def test_create_component_from_devfile_starter_project(self):

        print('Test case : should successfully create the devfile component and download the source in the context when used with --starter flag')

        # example devfile path
        source_devfile_path = os.path.join(os.path.dirname(__file__),
                                           '../examples/source/devfiles/java-openliberty/devfile.yaml')

        source_devfile_path = os.path.abspath(source_devfile_path)

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            starter_project = get_starter_project(source_devfile_path)
            print('starter project extracted from devfile:', starter_project)
            subprocess.run(["odo", "create", "--devfile", source_devfile_path, "--starter", starter_project])

    def test_create_component_with_context_starter_project(self):

        print('Test case : should successfully create the devfile component and download the source in the context when used with --starter flag')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            context = os.path.join(tmp_workspace, 'newcontext')
            subprocess.run(["odo", "create", "nodejs", "--starter", "nodejs-starter", "--context", context])

            list_files: list[str] = [
                "package.json",
                "package-lock.json",
                "README.md",
                "devfile.yaml"
            ]

            assert check_files_exist(list_files, context)


    def test_create_component_with_starter_project_git_branch(self):

        print('Test case : should successfully create the component and download the source from the specified branch')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-with-branch.yaml')
            copy_example_devfile(source_devfile_path, tmp_workspace)
            subprocess.run(["odo", "create", "nodejs", "--starter", "nodejs-starter"])

            list_files: list[str] = [
                "package.json",
                "package-lock.json",
                "README.md",
                "devfile.yaml"
            ]
            assert check_files_exist(list_files, tmp_workspace)


    def test_create_component_with_devfile_flag(self):

        print('Test case : should successfully create the devfile component with --devfile points to a devfile')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile.yaml')
            copy_example_devfile(source_devfile_path, tmp_workspace)

            subprocess.run(["odo", "create", "nodejs", "--devfile", "./devfile.yaml"])


    def test_create_component_with_context_json_output(self):

        print('Test case : should successfully create the devfile component and download the source in the context when used with --starter flag and generate json output')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            context = os.path.join(tmp_workspace, 'newcontext')
            result = subprocess.run(["odo", "create", "nodejs", "--starter", "nodejs-starter", "--context", context, "-o", "json"],
                                    capture_output=True, text=True, check=True)

            list_components = [
                "Component",
                "nodejs",
                "Not Pushed"
            ]

            assert validate_json_format(result.stdout)
            assert match_all(result.stdout, list_components)


    def test_create_component_with_devfile_url_flag(self):

        print('Test case : should successfully create the devfile component with --devfile points to a devfile url')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            subprocess.run(
                ["odo", "create", "nodejs", "--devfile", "https://raw.githubusercontent.com/odo-devfiles/registry/master/devfiles/nodejs/devfile.yaml"])
