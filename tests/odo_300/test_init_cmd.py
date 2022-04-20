import sys
import tempfile
import jmespath
from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry_v300")
class TestInitCmd:

    CONTEXT = "test-context"
    COMPONENT = "acomponent"

    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])

    def test_init_with_devfile_flag(self):

        print("Test case : should download a devfile.yaml file and correctly set the component name in it")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile", "go"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))
            assert check_file_exist("devfile.yaml", tmp_workspace)

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == self.COMPONENT


    def test_init_with_devfile_flag_json_output(self):
        print("Test case : should download a devfile.yaml file and correctly set the component name in it")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            result_json = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile", "go", "-o", "json"],
                                    capture_output=True, text=True, check=True)

            assert check_file_exist("devfile.yaml", tmp_workspace)
            assert validate_json_format(result_json.stdout)

            dict = json.loads(result_json.stdout)

            assert contains(jmespath.search('devfilePath', dict), os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml')))
            assert contains(jmespath.search('devfileData.devfile.schemaVersion', dict), "2.0.0")
            assert jmespath.search('devfileData.supportedOdoFeatures.dev', dict)
            assert not jmespath.search('devfileData.supportedOdoFeatures.debug', dict)
            assert not jmespath.search('devfileData.supportedOdoFeatures.deploy', dict)
            assert contains(jmespath.search('managedBy', dict), "odo")


    def test_init_with_devfile_path(self):
        print("Test case : using --devfile-path flag with a local devfile")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = get_source_devfile_path('nodejs/devfile-registry.yaml')

            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert check_file_exist("devfile.yaml", tmp_workspace)
            assert contains(result.stdout,
                            "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == self.COMPONENT


    def test_init_with_devfile_path_url(self):
        print("Test case : using --devfile-path flag with a URL")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = "https://raw.githubusercontent.com/devfile/integration-tests/main/tests/examples/source/devfiles/nodejs/devfile.yaml"

            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert check_file_exist("devfile.yaml", tmp_workspace)
            assert contains(result.stdout,
                            "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == self.COMPONENT


    def test_init_with_devfile_registry(self):
        print("Test case : should successfully run odo init if specified registry is valid")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile", "go",
                                     "--devfile-registry", "DefaultDevfileRegistry"],
                                    capture_output=True, text=True, check=True)

            assert check_file_exist("devfile.yaml", tmp_workspace)
            assert contains(result.stdout,
                            "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == self.COMPONENT


    def test_init_with_starter(self):

        print("Test case : should pass and keep the devfile in starter")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-with-starter-with-devfile.yaml")
            
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))
            assert check_file_exist("devfile.yaml", tmp_workspace)

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == self.COMPONENT
            assert query_yaml(devfile_path, "metadata", "language", -1) == "nodejs"

            list_contents: list[str] = [
                "2.2.0", "outerloop-deploy", "deployk8s", "outerloop-build"
            ]
            assert match_strings_in_file(devfile_path, list_contents)


    def test_init_with_starter_subdir(self):

        print("Test case : should successfully extract the project in the specified subDir path")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("springboot/devfile-with-subDir.yaml")
            
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--starter", "springbootproject", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            list_files: list[str] = [
                "java/com/example/demo/DemoApplication.java",
                "resources/application.properties"
            ]
            assert check_files_exist(list_files, tmp_workspace)


    def test_init_with_starter_and_branch(self):

        print("Test case : should successfully run odo init for devfile with starter project from the specified branch")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-with-branch.yaml")
            
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            list_files: list[str] = [
                "package.json", "package-lock.json", "README.md", "devfile.yaml", "test"
            ]
            assert check_files_exist(list_files, tmp_workspace)


    def test_init_with_starter_and_tag(self):

        print("Test case : should successfully run odo init for devfile with starter project from the specified tag")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-with-tag.yaml")
            
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            list_files: list[str] = [
                "package.json", "package-lock.json", "README.md", "devfile.yaml", "app"
            ]
            assert check_files_exist(list_files, tmp_workspace)


    def test_init_with_sources(self):

        print("Test case : running odo init from a directory with sources")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            copy_example("nodejs/project", tmp_workspace, self.CONTEXT)
            os.chdir(self.CONTEXT)

            # should work without --starter flag
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile", "nodejs"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            # should not accept --starter flag
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile", "nodejs", "--starter", "nodejs-starter"],
                                    capture_output=True, text=True, check=False)

            assert contains(result.stderr, "a devfile already exists in the current directory")

