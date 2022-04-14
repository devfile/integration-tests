import sys
import tempfile
from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry_v300")
class TestInitCmd:

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

    def test_init_with_devfile_flag(self):

        print("Test case : should download a devfile.yaml file and correctly set the component name in it")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            compName = "acomponent"
            result = subprocess.run(["odo", "init", "--name", compName, "--devfile", "go"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))
            assert check_file_exist(tmp_workspace, "devfile.yaml")

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == compName


    def test_init_with_starter(self):

        print("Test case : should pass and keep the devfile in starter")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-with-starter-with-devfile.yaml')

            compName = "acomponent"
            result = subprocess.run(["odo", "init", "--name", compName, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))
            assert check_file_exist(tmp_workspace, "devfile.yaml")

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))
            assert query_yaml(devfile_path, "metadata", "name", -1) == compName
            assert query_yaml(devfile_path, "metadata", "language", -1) == "nodejs"

            list_contents: list[str] = [
                "2.2.0", "outerloop-deploy", "deployk8s", "outerloop-build"
            ]
            assert match_strings_in_file(devfile_path, list_contents)


    def test_init_with_starter_subdir(self):

        print("Test case : should successfully extract the project in the specified subDir path")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/springboot/devfile-with-subDir.yaml')
            compName = "acomponent"
            result = subprocess.run(["odo", "init", "--name", compName, "--starter", "springbootproject", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))

            list_files: list[str] = [
                "java/com/example/demo/DemoApplication.java",
                "resources/application.properties"
            ]
            assert check_files_exist(tmp_workspace, list_files)


    def test_init_with_starter_and_branch(self):

        print("Test case : should successfully run odo init for devfile with starter project from the specified branch")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-with-branch.yaml')
            compName = "acomponent"
            result = subprocess.run(["odo", "init", "--name", compName, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))

            list_files: list[str] = [
                "package.json", "package-lock.json", "README.md", "devfile.yaml", "test"
            ]
            assert check_files_exist(tmp_workspace, list_files)


    def test_init_with_starter_and_tag(self):

        print("Test case : should successfully run odo init for devfile with starter project from the specified tag")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-with-tag.yaml')
            compName = "acomponent"
            result = subprocess.run(["odo", "init", "--name", compName, "--starter", "nodejs-starter", "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))

            list_files: list[str] = [
                "package.json", "package-lock.json", "README.md", "devfile.yaml", "app"
            ]
            assert check_files_exist(tmp_workspace, list_files)


    def test_init_with_sources(self):

        print("Test case : running odo init from a directory with sources")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)
            copy_example("nodejs/project", tmp_workspace, self.CONTEXT)
            os.chdir(self.CONTEXT)

            compName = "acomponent"

            # should work without --starter flag
            result = subprocess.run(["odo", "init", "--name", compName, "--devfile", "nodejs"],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(compName))

            # should not accept --starter flag
            result = subprocess.run(["odo", "init", "--name", compName, "--devfile", "nodejs", "--starter", "nodejs-starter"],
                                    capture_output=True, text=True, check=False)

            assert contains(result.stderr, "a devfile already exists in the current directory")
