from utils.config import *
from utils.util import *

@pytest.mark.usefixtures("use_test_registry")
class TestCatalogCmd:

    def test_catalog_cmd_list_components(self):

        print('Test case : should list all supported devfile components')

        result = subprocess.run(["odo", "catalog", "list", "components"],
                                capture_output=True, text=True, check=True)

        print('odo catalog list components :', result.stdout)

        list_components = [
            "Odo Devfile Components",
            "NAME",
            "REGISTRY",
            "DefaultDevfileRegistry",
            "java-maven",
            "java-openliberty",
            "java-quarkus",
            "java-springboot",
            "java-vertx",
            "java-wildfly",
            "java-wildfly-bootable-jar",
            "nodejs",
            "python",
            "python-django",
        ]

        assert match_all(result.stdout, list_components)


    def test_catalog_cmd_list_components_json_format(self):

        print('Test case : should list devfile components in json format')

        result = subprocess.run(["odo", "catalog", "list", "components", "-o", "json"],
                                capture_output=True, text=True, check=True)

        list_components = [
            "odo.dev/v1alpha1",
            "java-maven",
            "java-openliberty",
            "java-quarkus",
            "java-springboot",
            "java-vertx",
            "java-wildfly",
            "java-wildfly-bootable-jar",
            "nodejs",
            "python",
            "python-django",
        ]

        assert validate_json_format(result.stdout)
        assert match_all(result.stdout, list_components)


    def test_catalog_cmd_list_describe_json_format(self):

        print('Test case : when executing catalog describe component with -o json, it should display a valid JSON')

        result = subprocess.run(["odo", "catalog", "describe", "component", "nodejs", "-o", "json"],
                                capture_output=True, text=True, check=True)

        list_descriptions: list[str] = [
            "Node.js Runtime",
            "Stack with Node.js",
        ]

        assert validate_json_format(result.stdout)
        assert match_all(result.stdout, list_descriptions)
