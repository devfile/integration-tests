import sys
import tempfile
from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry_v300")
class TestDeleteCmd:

    CONTEXT = "test-context"
    COMPONENT = "acomponent"
    PROJECT = "intg-test-project"

    tmp_project_name = None

    @classmethod
    def setup_class(cls):
        # Runs once per class
        cls.tmp_project_name = create_test_project()

    @classmethod
    def teardown_class(cls):
        '''Runs at end of class'''
        subprocess.run(["odo", "project", "delete", cls.tmp_project_name, "-f", "-w"])

    @pytest.mark.skip(reason="This testcase is deprecated as it's covered in Odo integration tests")
    def test_delete(self):
        print("Test case : the component is deployed in DEV mode and it is deleted using its name and namespace")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-deploy-with-multiple-resources.yaml")
            copy_example("nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"
            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            cmd_odo_dev = str('odo dev --random-ports')

            try:
                # starting dev mode with random ports should work
                cmd_proc = subprocess.Popen(cmd_odo_dev, shell=True, bufsize=-1)
                # sleep until odo dev is started
                time.sleep(3)
                result = subprocess.run(["minikube",  "kubectl", "--",  "get", "deployment", "-n", self.PROJECT],
                        capture_output=True, text=True, check=True)

                list_expected = [
                    self.COMPONENT + "-app",
                    "NAME",
                    "AVAILABLE",
                ]
                assert match_all(result.stdout, list_expected)

                result = subprocess.run(["odo", "delete", "component", "--name", self.COMPONENT, "--namespace", self.PROJECT, "-f"],
                                        capture_output=True, text=True, check=True)
                str_deleted = "The component \"{}\" is successfully deleted from namespace \"{}\"".format(self.COMPONENT, self.PROJECT)

                assert contains(result.stdout, str_deleted)
                cmd_proc.terminate()

            except Exception as e:
                raise e

    @pytest.mark.skip(reason="This testcase is deprecated as it's covered in Odo integration tests")
    def test_delete_from_other_directory(self):
        print("Test case : the component is deployed in DEV mode and it is deleted using its name and namespace from another directory")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-deploy-with-multiple-resources.yaml")
            copy_example("nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"
            os.chdir(self.CONTEXT)

            # should work without --starter flag
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            cmd_odo_dev = str('odo dev --random-ports')

            try:
                # starting dev mode with random ports should work
                cmd_proc = subprocess.Popen(cmd_odo_dev, shell=True, bufsize=-1)
                # sleep while odo dev is started
                time.sleep(3)
                result = subprocess.run(["minikube",  "kubectl", "--",  "get", "deployment", "-n", self.PROJECT],
                        capture_output=True, text=True, check=True)

                list_expected = [
                    self.COMPONENT + "-app",
                    "NAME",
                    "AVAILABLE",
                ]
                assert match_all(result.stdout, list_expected)

                os.chdir(tmp_workspace)
                os.mkdir("test-dir")
                os.chdir("test-dir")
                result = subprocess.run(["odo", "delete", "component", "--name", self.COMPONENT, "--namespace", self.PROJECT, "-f"],
                                        capture_output=True, text=True, check=True)
                str_deleted = "The component \"{}\" is successfully deleted from namespace \"{}\"".format(self.COMPONENT, self.PROJECT)

                assert contains(result.stdout, str_deleted)

                result = subprocess.run(["minikube",  "kubectl", "--",  "get", "deployment", "-n", self.PROJECT],
                                        capture_output=True, text=True, check=True)

                assert not contains(result.stdout, self.COMPONENT + "-app")
                assert contains(result.stderr, "No resources found in " + self.PROJECT + " namespace.")

                cmd_proc.terminate()
            except Exception as e:
                raise e

    @pytest.mark.skip(reason="This testcase is deprecated as it's covered in Odo integration tests")
    def test_delete_with_devfile_present(self):
        print("Test case : the component is deployed in DEV mode and it is deleted when devfile present")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # get test devfile path
            source_devfile_path = get_source_devfile_path("nodejs/devfile-deploy-with-multiple-resources.yaml")
            copy_example("nodejs/project", tmp_workspace, self.CONTEXT)

            os.environ['PODMAN_CMD'] = "echo"
            os.chdir(self.CONTEXT)

            # should work without --starter flag
            result = subprocess.run(["odo", "init", "--name", self.COMPONENT, "--devfile-path", source_devfile_path],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "Your new component '{}' is ready in the current directory.".format(self.COMPONENT))

            cmd_odo_dev = str('odo dev --random-ports')

            try:
                # starting dev mode with random ports should work
                cmd_proc = subprocess.Popen(cmd_odo_dev, shell=True, bufsize=-1)
                # sleep while odo dev is started
                time.sleep(3)
                result = subprocess.run(["minikube",  "kubectl", "--",  "get", "deployment", "-n", self.PROJECT],
                        capture_output=True, text=True, check=True)

                list_expected = [
                    self.COMPONENT + "-app",
                    "NAME",
                    "AVAILABLE",
                ]
                assert match_all(result.stdout, list_expected)

                str_deleted = "The component \"{}\" is successfully deleted from namespace \"{}\"".format(self.COMPONENT, self.PROJECT)
                result = subprocess.run(["odo", "delete", "component", "-f"],
                                        capture_output=True, text=True, check=True)

                assert contains(result.stdout, str_deleted)

                result = subprocess.run(["minikube",  "kubectl", "--",  "get", "deployment", "-n", self.PROJECT],
                                        capture_output=True, text=True, check=True)

                assert not contains(result.stdout, self.COMPONENT + "-app")
                assert contains(result.stderr, "No resources found in " + self.PROJECT + " namespace.")

                # run delete command again after deletion
                result = subprocess.run(["odo", "delete", "component", "-f"],
                    capture_output=True, text=True, check=True)
                str_no_resource_found = "No resource found for component \"{}\" in namespace \"{}\"".format(self.COMPONENT, self.PROJECT)
                assert contains(result.stdout, str_no_resource_found)

                cmd_proc.terminate()
            except Exception as e:
                raise e
