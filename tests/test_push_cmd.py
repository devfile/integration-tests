import sys
import tempfile

from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestPushCmd:

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

    @pytest.mark.skipif(sys.platform == "linux", reason="Fixme: kubectl call does not return correct value on linux")
    def test_push_with_devfile(self):

        print("Test case : checks that odo push works with a devfile")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-registry.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "push", "--project", self.tmp_project_name],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Changes successfully pushed to component")

            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "deployment",
                 "-o", "jsonpath='{.items[0].spec.template.spec.containers[0].ports[?(@.name=='3000-tcp')].containerPort}'"],
                capture_output=True, text=True, check=True)
            assert contains(result.stdout, "3000")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")


    def test_push_with_json_output(self):

        print("Test case : check that odo push is executed with json output")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "push", "-o", "json", "--project", self.tmp_project_name],
                                    capture_output=True, text=True, check=True)

            assert contains(result.stdout, "devFileCommandExecutionComplete")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")


    @pytest.mark.skipif(sys.platform == "linux", reason="Fixme: kubectl call does not return correct value on linux")
    def test_push_with_env_variable(self):

        print("Test case : should check if the env variable has a correct value after push")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-variables.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "push", "--project", self.tmp_project_name],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Changes successfully pushed to component")

            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "deployment",
                 "-o", "jsonpath='{.items[0].spec.template.spec.containers[0].env[?(@.name == 'FOO')].value}'"],
                capture_output=True, text=True, check=True)
            assert contains(result.stdout, "bar")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")

    @pytest.mark.skipif(sys.platform == "linux", reason="Fixme: kubectl call does not return correct value on linux")
    def test_push_with_sourcemappings(self):

        print("Test case : when devfile has sourcemappings and doing odo push it should sync files to the correct location")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfileSourceMapping.yaml')

            copy_and_create(source_devfile_path, "nodejs/project", tmp_workspace, self.CONTEXT)

            os.chdir(self.CONTEXT)

            result = subprocess.run(["odo", "push", "--project", self.tmp_project_name],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Changes successfully pushed to component")

            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "deployment",
                 "-o", "jsonpath='{.items[0].spec.template.spec.containers[0].env[?(@.name == 'PROJECT_SOURCE')].value}'"],
                capture_output=True, text=True, check=True)
            assert contains(result.stdout, "/test")

            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "pods",
                 "-o", "jsonpath='{.items[*].metadata.name}'"],
                capture_output=True, text=True, check=True)
            pod_name = result.stdout.strip("\'")

            result = subprocess.run(
                ["minikube", "kubectl", "--", "exec", pod_name,
                 "--", "stat", "/test/server.js"],
                capture_output=True, text=True, check=True)
            assert contains(result.stdout, "File: /test/server.js")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")


    # project and clonePath is present in devfile and doing odo push
    @pytest.mark.skipif(sys.platform == "linux", reason="Fixme: kubectl call does not return correct value on linux")
    def test_push_with_clonepath(self):

        print("Test case : should sync to the correct dir in container")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-with-projects.yaml')
            copy_example_devfile(source_devfile_path, tmp_workspace)

            result = subprocess.run(["odo", "push", "--v", "5"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Changes successfully pushed to component")

            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "deployment",
                 "-o", "jsonpath='{.items[0].spec.template.spec.containers[0].env[?(@.name == 'PROJECT_SOURCE')].value}'"],
                capture_output=True, text=True, check=True)
            assert result.stdout.strip("\'") == "/apps/webapp"

            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "pods",
                 "-o", "jsonpath='{.items[*].metadata.name}'"],
                capture_output=True, text=True, check=True)
            pod_name = result.stdout.strip("\'")

            result = subprocess.run(
                ["minikube", "kubectl", "--", "exec", pod_name,
                 "--", "ls", "/apps/webapp"],
                capture_output=True, text=True, check=True)

            # source code is synced to $PROJECTS_ROOT/[clonePath] where $PROJECTS_ROOT is '/projects' by default.
            # When 'sourceMapping' and 'clonePath' are set in 'devfile-with-projects.yaml', the
            # source code should be synced to '/apps/webapp' (i.e. sourceMapping: /apps, sourceMapping: /apps)
            assert contains(result.stdout, "package-lock.json")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")


    # devfile project field is present with no clonePath in devfile and doing odo push
    @pytest.mark.skipif(sys.platform == "linux", reason="Fixme: kubectl call does not return correct value on linux")
    def test_push_with_project(self):

        print("Test case : should sync to the correct dir in container")

        with tempfile.TemporaryDirectory() as tmp_workspace:
            os.chdir(tmp_workspace)

            # example devfile path
            source_devfile_path = os.path.join(os.path.dirname(__file__),
                                               'examples/source/devfiles/nodejs/devfile-with-projects.yaml')
            copy_example_devfile(source_devfile_path, tmp_workspace)
            replace_string_in_a_file("devfile.yaml", "clonePath: webapp/", "# clonePath: webapp/")
            result = subprocess.run(["odo", "push"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Changes successfully pushed to component")

            # MacOS: reuse the existing kubectl inside minikube
            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "deployment",
                 "-o", "jsonpath='{.items[0].spec.template.spec.containers[0].env[?(@.name == 'PROJECT_SOURCE')].value}'"],
                capture_output=True, text=True, check=True)
            assert result.stdout.strip("\'") == "/apps/nodeshift"

            result = subprocess.run(
                ["minikube", "kubectl", "--", "get", "pods",
                 "-o", "jsonpath='{.items[*].metadata.name}'"],
                capture_output=True, text=True, check=True)
            pod_name = result.stdout.strip("\'")

            result = subprocess.run(
                ["minikube", "kubectl", "--", "exec", pod_name,
                 "--", "ls", "/apps/nodeshift"],
                capture_output=True, text=True, check=True)

            # source code is synced to $PROJECTS_ROOT/[project_name] where $PROJECTS_ROOT is '/projects' by default.
            # When 'sourceMapping' is set to '/apps' and no 'clonePath'in 'devfile-with-projects.yaml', the
            # source code should be synced to '/apps/nodeshift' (i.e. sourceMapping: /apps, projects:name is 'nodeshift')
            assert contains(result.stdout, "package-lock.json")

            result = subprocess.run(["odo", "delete", "-f", "-w"],
                                    capture_output=True, text=True, check=True)
            assert contains(result.stdout, "Successfully deleted component")
