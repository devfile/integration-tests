import tempfile

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestConfigCmd:
    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])

    def test_config_view(self):

        print('Test case : should view all default parameters')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-registry.yaml')
            shutil.copyfile(source_devfile_path, devfile_path)
            subprocess.check_call(["odo", "create", "nodejs", "--devfile", "./devfile.yaml"])

            result = subprocess.run(["odo", "config", "view"],
                                    capture_output=True, text=True, check=True)

            print('odo config view :', result.stdout)

            list_components = [
                "nodejs",
                "Ports",
                "Memory",
            ]

            assert match_all(result.stdout, list_components)


    def test_config_set(self):

        TESTNAME = "testname"
        MEMORY = "500Mi"
        PORT = "8888"

        print('Test case : hould successfully set the parameters')

        with tempfile.TemporaryDirectory() as tmp_workspace:
            print('created temporary workspace', tmp_workspace)
            os.chdir(tmp_workspace)

            devfile_path = os.path.abspath(os.path.join(tmp_workspace, 'devfile.yaml'))

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               '../examples/source/devfiles/nodejs/devfile-registry.yaml')
            shutil.copyfile(source_devfile_path, devfile_path)
            subprocess.check_call(["odo", "create", "nodejs", "--devfile", "./devfile.yaml"])

            subprocess.check_call(["odo", "config", "set", "Name", TESTNAME, "-f"])
            subprocess.check_call(["odo", "config", "set", "Memory", MEMORY, "-f"])
            subprocess.check_call(["odo", "config", "set", "Ports", PORT, "-f"])

            result = subprocess.run(["odo", "config", "view"],
                                    capture_output=True, text=True, check=True)

            print('odo config view :', result.stdout)

            list_components = [
                TESTNAME,
                MEMORY,
                PORT,
            ]

            assert match_all(result.stdout, list_components)

            subprocess.check_call(["odo", "config", "unset", "Ports", "-f"])
            result = subprocess.run(["odo", "config", "view"],
                                    capture_output=True, text=True, check=True)

            list_components = [
                PORT,
            ]
            assert unmatch_all(result.stdout, list_components)
