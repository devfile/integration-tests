import json
import os
import random
import shutil
import subprocess
import string
import time
import yaml
from pathlib import Path

def get_starter_project(devfile_path: str):
    return query_yaml(devfile_path, "starterProjects", 0, "name")

# True if str contains s
def contains(str, s):
    if str is None:
        return False
    return str.__contains__(s)

# match_all ensures all expected strings are found in input
def match_all(input, list_expected):
    for item in list_expected:
        if not contains(input, item):
            return False
    return True

# unmatch_all ensures input doesn't contain any string in the list
def unmatch_all(input, list_not_expected):
    for item in list_not_expected:
        if contains(input, item):
            return False
    return True

# This function loads the given devfile
def get_devfile(devfile_path = 'devfile.yaml'):
    with open(devfile_path, 'r') as file:
        devfile_data = yaml.safe_load(file)
    return devfile_data

# This function loads the given schema available
def get_schema_json():
    with open('devfile.json', 'r') as file:
        schema = json.load(file)
    return schema

def validate_json():
    execute_api_schema = get_schema_json()
    json_data = get_devfile()

    try:
        validate(instance=json_data, schema=execute_api_schema)
    except jsonschema.exceptions.ValidationError as err:
        print(err)
        err = "Given JSON data is Invalid"
        return False, err

    message = "Given JSON data is Valid"
    return True, message

def validate_json_format(json_data):
    try:
        json.loads(json_data)
    except ValueError as err:
        return False
    return True

def set_default_devfile_registry():
    # Use staging devfile repository as the default repository for tests
    registryName = "DefaultDevfileRegistry"

    # Use staging OCI - based registry
    default_devfile_registry = "https://registry.stage.devfile.io"

    result = subprocess.run(["odo", "registry", "add", registryName, default_devfile_registry])

# generate a randome string
def random_string(size=8, chars=string.ascii_lowercase + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))

def query_yaml(devfile_path: str, param_1, param_2, param_3):
    try:
        devfile_data = get_devfile(devfile_path)
        print('devfile_data content:', devfile_data)

        if param_2 == -1:
            return devfile_data[param_1]
        elif param_3 == -1:
            return devfile_data[param_1][param_2]
        else:
            return devfile_data[param_1][param_2][param_3]
    except yaml.YAMLError as e:
        print(e)

# check if files exist in the context
def check_files_exist(context, list_files):

    for expected in list_files:
        path_to_file = os.path.join(context, expected)
        if os.path.exists(path_to_file):
            continue
        return False
    return True

# wait until the file exists
# def wait_for_file(source_dir, filename, timeout = 20):
#
#     file_path = os.path.normpath(os.path.join(source_dir, filename))
#     attempts = 0
#
#     while attempts < timeout:
#         # Check if the file exists.
#         if os.path.isfile(file_path):
#             return
#         # Wait 1 second before trying again.
#         time.sleep(1)
#         attempts += 1

# replace_string_in_a_file replaces old_string with new_string in text file
def replace_string_in_a_file(filename, old_string, new_string):
    fin = open(filename, "rt")
    data = fin.read()
    data = data.replace(old_string, new_string)
    fin.close()
    fin = open(filename, "wt")
    fin.write(data)
    fin.close()

# Copy example files to context directory
def copy_example(example_name, workspace_dir, context_dir = '.'):
    example_path = os.path.join(os.path.dirname(__file__), '../',
                                       'examples/source/devfiles', example_name)
    target_path = os.path.join(workspace_dir, context_dir)
    ''' context_dir shouldn't exist before calling copytree'''
    shutil.copytree(example_path, target_path)

# Copy example devfile.yaml to context directory
def copy_example_devfile(source_devfile, workspace_dir, context_dir = '.'):
    devfile_path = os.path.abspath(os.path.join(workspace_dir, context_dir, 'devfile.yaml'))
    shutil.copyfile(source_devfile, devfile_path)

# Copy example files and devfile.yaml to context directory and create component
def copy_and_create(source_devfile, example_name, workspace_dir, context_dir = '.'):
    copy_example(example_name, workspace_dir, context_dir)
    copy_example_devfile(source_devfile, workspace_dir, context_dir)
    subprocess.run(["odo", "create", "--context", os.path.join(workspace_dir, context_dir)])
