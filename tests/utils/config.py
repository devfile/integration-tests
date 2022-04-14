import pytest
import subprocess

from subprocess import CompletedProcess
from typing import Union, Any
from utils.util import *

@pytest.fixture(scope="session")
def use_test_registry():

    # Use staging OCI-based registry for tests to avoid overload
    test_registry_name = "TestDevfileRegistry"
    test_registry_url = "https://registry.stage.devfile.io"

    print('Check if ', test_registry_name, 'exists and update URL if needed.')

    result = subprocess.run(["odo", "registry", "list"],
                   capture_output=True, text=True, check=True)

    if contains(result.stdout, test_registry_name):
        print('Updating', test_registry_name, '...')

        subprocess.run(["odo", "registry", "update", test_registry_name, test_registry_url, '-f'],
                   capture_output=True, text=True, check=True)
    else:
        print('Creating', test_registry_name, '...')
        subprocess.run(["odo", "registry", "add", test_registry_name, test_registry_url],
                   capture_output=True, text=True, check=True)
    print('Using', test_registry_name, ':', test_registry_url)

    yield test_registry_name


@pytest.fixture(scope="session")
def use_test_registry_v300():

    # Use staging OCI-based registry for tests to avoid overload
    test_registry_name = "TestDevfileRegistry"
    test_registry_url = "https://registry.stage.devfile.io"

    print('Check if', test_registry_name, 'exists and update URL if needed.')

    result = subprocess.run(["odo", "preference", "registry", "list"],
                            capture_output=True, text=True, check=True)

    if contains(result.stdout, test_registry_name):
        print('Updating', test_registry_name, '...')

        subprocess.run(["odo", "preference", "registry", "update", test_registry_name, test_registry_url, '-f'],
                       capture_output=True, text=True, check=True)
    else:
        print('Creating', test_registry_name, '...')
        subprocess.run(["odo", "preference", "registry", "add", test_registry_name, test_registry_url],
                       capture_output=True, text=True, check=True)

    print('Using', test_registry_name, ':', test_registry_url)

    yield test_registry_name


def create_test_project():

    tmp_project_name = "intg-test-project"
    subprocess.run(["odo", "project", "create", tmp_project_name, "-w"])
    return tmp_project_name
