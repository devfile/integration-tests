# app command
odo devfile app command tests when the user creates and pushes two new devfile components in different apps when the user creates and pushes a third devfile component should list, describe and delete the app properly with json output
odo devfile app command tests when the user creates and pushes two new devfile components in different apps when the user creates and pushes a third s2i component on a openshift cluster should list, describe and delete the app properly with json output

# catalog command
odo devfile catalog command tests When checking catalog for installed services should succeed
odo devfile catalog command tests When executing catalog describe component with -o json should display a valid JSON
odo devfile catalog command tests When executing catalog describe component with a component name with multiple components should print multiple devfiles from different registries
odo devfile catalog command tests When executing catalog list components should list all supported devfile components
odo devfile catalog command tests When executing catalog list components should list components successfully even with an invalid kubeconfig path or path points to existing directory
odo devfile catalog command tests When executing catalog list components with -o json flag should list devfile components in json format
odo devfile catalog command tests When executing catalog list components with registry that is not set up properly should list components from valid registry

# config command
odo devfile config command tests When executing config set and unset Should fail to set and unset an invalid parameter
odo devfile config command tests When executing config set and unset Should successfully set and unset the parameters
odo devfile config command tests When executing config view Should view all default parameters

# create command
odo devfile create command tests When executing odo create using --starter with a devfile component that contains no projects should fail with please run &#39;no starter project found in devfile.&#39;
odo devfile create command tests When executing odo create with --s2i flag should fail to create the component specified with valid project and download the source
odo devfile create command tests When executing odo create with --s2i flag should fail to create the devfile component which doesn&#39;t have an s2i component of same name
odo devfile create command tests When executing odo create with --s2i flag should fail to create the devfile component with --registry specified
odo devfile create command tests When executing odo create with --s2i flag should fail to create the devfile component with valid file system path
odo devfile create command tests When executing odo create with an invalid project specified in --starter should fail with please run &#39;The project: invalid-project-name specified in --starter does not exist&#39;
odo devfile create command tests When executing odo create with component with no devBuild command should successfully create the devfile component and remove a dangling env file
odo devfile create command tests When executing odo create with devfile component and --starter flag should successfully create the component and download the source
odo devfile create command tests When executing odo create with devfile component and --starter flag should successfully create the component specified with valid project and download the source
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create and push the devfile component and show json output for working cluster
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component and download the source when used with --starter flag
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component and show json output for a unreachable cluster
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component and show json output for non connected cluster
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component and show json output for working cluster
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component in the context
odo devfile create command tests When executing odo create with devfile component type argument and --context flag should successfully create the devfile component with auto generated name
odo devfile create command tests When executing odo create with devfile component type argument and --project flag should successfully create the devfile component
odo devfile create command tests When executing odo create with devfile component type argument and --registry flag should fail to create the devfile component if specified registry is invalid
odo devfile create command tests When executing odo create with devfile component type argument and --registry flag should successfully create the devfile component if specified registry is valid
odo devfile create command tests When executing odo create with devfile component type argument should fail to create the devfile component with invalid component type
odo devfile create command tests When executing odo create with devfile component type argument should successfully create the devfile component with valid component name
odo devfile create command tests When executing odo create with devfile component, --starter flag and subDir has a valid value should only extract the specified path in the subDir field
odo devfile create command tests When executing odo create with existing devfile Testing Create for OpenShift specific scenarios should fail when we create the devfile or s2i component multiple times
odo devfile create command tests When executing odo create with existing devfile When devfile exists in user&#39;s working directory should fail to create the devfile component with --devfile points to different devfile
odo devfile create command tests When executing odo create with existing devfile When devfile exists in user&#39;s working directory should fail to create the devfile component with more than 1 arguments are passed in
odo devfile create command tests When executing odo create with existing devfile When devfile exists in user&#39;s working directory should fail when we create the devfile component multiple times
odo devfile create command tests When executing odo create with existing devfile When devfile exists in user&#39;s working directory should successfully create the devfile component
odo devfile create command tests When executing odo create with existing devfile When devfile exists in user&#39;s working directory should successfully create the devfile component with --devfile points to the same devfile
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should fail to create the devfile component with --registry specified
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should fail to create the devfile component with invalid URL path
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should fail to create the devfile component with invalid file system path
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should fail to create the devfile component with more than 1 arguments are passed in
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should successfully create the devfile component with valid file system path
odo devfile create command tests When executing odo create with existing devfile When devfile exists not in user&#39;s working directory and user specify the devfile path via --devfile should successfully create the devfile component with valid specifies URL path
odo devfile create command tests When executing odo create with git tag or git branch specified in starter project should successfully create the component and download the source from the specified branch
odo devfile create command tests When executing odo create with git tag or git branch specified in starter project should successfully create the component and download the source from the specified tag
odo devfile create command tests checks that odo push works with a devfile with now flag

# debug command
odo devfile debug command tests odo debug info should work on a odo component should start a debug session and run debug info on a closed debug session
odo devfile debug command tests odo debug info should work on a odo component should start a debug session and run debug info on a running debug session
odo devfile debug command tests odo debug on a nodejs:latest component check that machine output debug information works
odo devfile debug command tests odo debug on a nodejs:latest component should error out on devfile flag
odo devfile debug command tests odo debug on a nodejs:latest component should expect a ws connection when tried to connect on default debug port locally
odo devfile debug command tests when the debug command throws an error during push should wait and error out with some log
odo devfile delete command tests when a component is created devfile has preStop events when component is pushed should execute the preStop events

# delete command
odo devfile delete command tests when a component is created should delete the component
odo devfile delete command tests when a component is created should throw an error on an invalid delete command
odo devfile delete command tests when a component is created when the component has resources attached to it when the component is pushed should delete the component and its owned resources
odo devfile delete command tests when a component is created when the component has resources attached to it when the component is pushed should delete the component and its owned resources with --wait flag
odo devfile delete command tests when a component is created when the component is pushed deleting a component from other component directory when the second component is created when the second component is pushed should delete all the config files and component of the context directory with --all flag
odo devfile delete command tests when a component is created when the component is pushed deleting a component from other component directory when the second component is created when the second component is pushed should delete the context directory&#39;s component
odo devfile delete command tests when a component is created when the component is pushed deleting a component from other component directory when the second component is created when the second component is pushed should delete with the component name
odo devfile delete command tests when a component is created when the component is pushed should delete the component, env, odo folders and odo-index-file.json with --all flag
odo devfile delete command tests when component is created from an existing devfile present in its directory should not delete the devfile
odo devfile delete command tests when component is created with --devfile flag should successfully delete devfile
odo devfile delete command tests when the component is created in a non-existent project should let the user delete the local config files with -a flag
odo devfile delete command tests when the component is created in a non-existent project when deleting outside a component directory should let the user delete the local config files with --context flag

# describe command
odo devfile describe command tests When executing odo describe should describe the component when it is not pushed
odo devfile describe command tests When executing odo describe should describe the component when it is pushed
odo devfile describe command tests when running odo describe for machine readable output should show json output for non connected cluster
odo devfile describe command tests when running odo describe for machine readable output should show json output for working cluster

# env command
odo devfile env command tests When executing env set and unset Should fail to set and unset an invalid parameter
odo devfile env command tests When executing env set and unset Should successfully set and unset the parameters
odo devfile env command tests When executing env view Should view all default parameters

# exec command
odo devfile exec command tests When devfile exec command is executed should error out when a component is not present or when a devfile flag is used
odo devfile exec command tests When devfile exec command is executed should error out when a invalid command is given by the user
odo devfile exec command tests When devfile exec command is executed should error out when no command is given by the user
odo devfile exec command tests When devfile exec command is executed should execute the given command successfully in the container

# log command
odo devfile log command tests Verify odo log for devfile works should error out if component does not exist
odo devfile log command tests Verify odo log for devfile works should log debug command output
odo devfile log command tests Verify odo log for devfile works should log run command output and fail for debug command


# push command
odo devfile push command tests Handle devfiles with parent should handle a devfile with a parent and add a extra command
odo devfile push command tests Handle devfiles with parent should handle a devfile with a parent and override a composite command
odo devfile push command tests Handle devfiles with parent should handle a multi layer parent
odo devfile push command tests Handle devfiles with parent should handle a parent and override/append it&#39;s envs
odo devfile push command tests Pushing devfile without an .odo folder should be able to push based on metadata.name in devfile WITH a dash in the name
odo devfile push command tests Pushing devfile without an .odo folder should be able to push based on name passed
odo devfile push command tests Pushing devfile without an .odo folder should error out on devfile flag
odo devfile push command tests Testing Push for Kubernetes specific scenarios should push successfully project value is default
odo devfile push command tests Testing Push for OpenShift specific scenarios throw an error when the project value is default
odo devfile push command tests Testing Push with remote attributes should push only the mentioned files at the appropriate remote destination
odo devfile push command tests Verify devfile push works Ensure that push -f correctly removes local deleted files from the remote target sync folder
odo devfile push command tests Verify devfile push works checks that odo push with -o json displays machine readable JSON event output
odo devfile push command tests Verify devfile push works checks that odo push works outside of the context directory
odo devfile push command tests Verify devfile push works checks that odo push works with a devfile
odo devfile push command tests Verify devfile push works checks that odo push works with a devfile with composite commands
odo devfile push command tests Verify devfile push works checks that odo push works with a devfile with nested composite commands
odo devfile push command tests Verify devfile push works checks that odo push works with a devfile with parallel composite commands
odo devfile push command tests Verify devfile push works checks that odo push works with a devfile with sourcemapping set
odo devfile push command tests Verify devfile push works should be able to create a file, push, delete, then push again propagating the deletions
odo devfile push command tests Verify devfile push works should be able to handle a missing build command group
odo devfile push command tests Verify devfile push works should be able to push using the custom commands
odo devfile push command tests Verify devfile push works should build when no changes are detected in the directory and force flag is enabled
odo devfile push command tests Verify devfile push works should create pvc and reuse if it shares the same devfile volume name
odo devfile push command tests Verify devfile push works should delete the files from the container if its removed locally
odo devfile push command tests Verify devfile push works should err out on an event command not mapping to a devfile container component
odo devfile push command tests Verify devfile push works should err out on an event composite command mentioning an invalid child command
odo devfile push command tests Verify devfile push works should err out on an event not mentioned in the devfile commands
odo devfile push command tests Verify devfile push works should error out if the devfile has an invalid command group
odo devfile push command tests Verify devfile push works should error out on a missing run command group
odo devfile push command tests Verify devfile push works should error out on a wrong custom commands
odo devfile push command tests Verify devfile push works should error out on multiple or no default commands
odo devfile push command tests Verify devfile push works should execute PostStart commands if present and not execute when component already exists
odo devfile push command tests Verify devfile push works should execute PreStart commands if present during pod startup
odo devfile push command tests Verify devfile push works should execute commands with flags if the command has no group kind
odo devfile push command tests Verify devfile push works should execute the default build and run command groups if present
odo devfile push command tests Verify devfile push works should have no errors when no endpoints within the devfile, should create a service when devfile has endpoints
odo devfile push command tests Verify devfile push works should not build when no changes are detected in the directory and build when a file change is detected
odo devfile push command tests Verify devfile push works should not restart the application if it is hot reload capable
odo devfile push command tests Verify devfile push works should restart the application if it is not hot reload capable
odo devfile push command tests Verify devfile push works should restart the application if run mode is changed, regardless of hotReloadCapable value
odo devfile push command tests Verify devfile push works should run odo push successfully after odo push --debug
odo devfile push command tests Verify devfile push works should successfully create the devfile component with valid component name
odo devfile push command tests Verify devfile push works should throw a validation error for composite command indirectly referencing itself
odo devfile push command tests Verify devfile push works should throw a validation error for composite command referencing non-existent commands
odo devfile push command tests Verify devfile push works should throw a validation error for composite command that has invalid exec subcommand
odo devfile push command tests Verify devfile push works should throw a validation error for composite run commands
odo devfile push command tests Verify devfile volume components work should error out if a container component has volume mount that does not refer a valid volume component
odo devfile push command tests Verify devfile volume components work should error out when a wrong volume size is used
odo devfile push command tests Verify devfile volume components work should error out when duplicate volume components exist
odo devfile push command tests Verify devfile volume components work should successfully use the volume components in container components
odo devfile push command tests Verify devfile volume components work should throw a validation error for v1 devfiles
odo devfile push command tests Verify files are correctly synced Should ensure that files are correctly synced on pod redeploy, with force push specified
odo devfile push command tests Verify files are correctly synced Should ensure that files are correctly synced on pod redeploy, without force push specified
odo devfile push command tests Verify source code sync location Should sync to the correct dir in container if multiple project is present
odo devfile push command tests Verify source code sync location Should sync to the correct dir in container if no project is present
odo devfile push command tests Verify source code sync location Should sync to the correct dir in container if project and clonePath is present
odo devfile push command tests Verify source code sync location Should sync to the correct dir in container if project present
odo devfile push command tests exec commands with environment variables Should be able to exec command with environment variable with spaces
odo devfile push command tests exec commands with environment variables Should be able to exec command with multiple environment variables
odo devfile push command tests exec commands with environment variables Should be able to exec command with single environment variable
odo devfile push command tests push with listing the devfile component checks components in a specific app and all apps
odo devfile push command tests push with listing the devfile component checks devfile and s2i components together
odo devfile push command tests when .gitignore file exists checks that .odo/env exists in gitignore
odo devfile push command tests when the run command throws an error should wait and error out with some log

# registry command
odo devfile registry command tests When executing registry commands with the registry is not present Should fail to delete the registry
odo devfile registry command tests When executing registry commands with the registry is not present Should fail to update the registry
odo devfile registry command tests When executing registry commands with the registry is not present Should successfully add the registry
odo devfile registry command tests When executing registry commands with the registry is present Should fail to add the registry
odo devfile registry command tests When executing registry commands with the registry is present Should successfully delete the registry
odo devfile registry command tests When executing registry commands with the registry is present Should successfully update the registry
odo devfile registry command tests When executing registry list Should fail with an error with no registries
odo devfile registry command tests When executing registry list Should list all default registries
odo devfile registry command tests When executing registry list Should list all default registries with json
odo devfile registry command tests when working with git based registries should not show deprecation warning if non-git-based registry is used
odo devfile registry command tests when working with git based registries should show deprecation warning when the git based registry is used

# status command
odo devfile status command tests Verify URL status is correctly reported Verify that odo component status detects the URL status: Ingress Nonsecure
odo devfile status command tests Verify URL status is correctly reported Verify that odo component status detects the URL status: Ingress Secure
odo devfile status command tests Verify URL status is correctly reported Verify that odo component status detects the URL status: Route Nonsecure
odo devfile status command tests Verify URL status is correctly reported Verify that odo component status detects the URL status: Route Secure
odo devfile status command tests Verify devfile status works Verify that odo component status correctly detects component Kubernetes pods
odo devfile status command tests Verify devfile status works Verify that odo component status correctly reports supervisord status

# storage command
odo devfile storage command tests When devfile storage commands are invalid should error if same path is provided again
odo devfile storage command tests When devfile storage commands are invalid should error if same storage name is provided again
odo devfile storage command tests When devfile storage commands are invalid should throw error if no storage is present
odo devfile storage command tests When devfile storage create command is executed should create a storage when storage is not provided
odo devfile storage command tests When devfile storage create command is executed should create a storage with default size when --size is not provided
odo devfile storage command tests When devfile storage create command is executed should create and output in json format
odo devfile storage command tests When devfile storage create command is executed should create storage and attach to specified container successfully and list it correctly
odo devfile storage command tests When devfile storage create command is executed should create the storage and mount it on the container
odo devfile storage command tests When devfile storage list command is executed should list output in json format
odo devfile storage command tests When devfile storage list command is executed should list the storage with the proper states
odo devfile storage command tests When devfile storage list command is executed should list the storage with the proper states and container names
odo devfile storage command tests When ephemeral is not set in preference.yaml should not create a pvc to store source code (default is ephemeral=false)
odo devfile storage command tests When ephemeral is set to false in preference.yaml should create a pvc to store source code
odo devfile storage command tests When ephemeral is set to true in preference.yaml should not create a pvc to store source code

# test command
odo devfile test command tests Should run test command successfully Should run composite test command successfully
odo devfile test command tests Should run test command successfully Should run test command successfully with only one default specified
odo devfile test command tests Should run test command successfully Should run test command successfully with test-command specified
odo devfile test command tests Should show proper errors should error out on devfile flag
odo devfile test command tests Should show proper errors should show error if component is not pushed; should error out if a non-existent command or a command from wrong group is specified
odo devfile test command tests Should show proper errors should show error if devfile has multiple default test command
odo devfile test command tests Should show proper errors should show error if devfile has no default test command
odo devfile test command tests Should show proper errors should show error if no test group is defined

# url command
odo devfile url command tests Creating urls create and delete with now flag should pass
odo devfile url command tests Creating urls should be able to push again twice after creating and deleting a url
odo devfile url command tests Creating urls should create URL with path defined in Endpoint
odo devfile url command tests Creating urls should create URLs under different container names
odo devfile url command tests Creating urls should create a URL without port flag if only one port exposed in devfile
odo devfile url command tests Creating urls should create a secure URL
odo devfile url command tests Creating urls should error out on devfile flag
odo devfile url command tests Creating urls should not allow creating an endpoint with same name
odo devfile url command tests Creating urls should not allow creating an invalid host
odo devfile url command tests Creating urls should not allow creating under an invalid container
odo devfile url command tests Creating urls should not allow using tls secret if url is not secure
odo devfile url command tests Creating urls should not create URLs under different container names with same port number
odo devfile url command tests Creating urls should report multiple issues when it&#39;s the case
odo devfile url command tests Listing urls should be able to list ingress url in machine readable json format
odo devfile url command tests Listing urls should list ingress url with appropriate state
odo devfile url command tests Listing urls should list url after push using context
odo devfile url command tests Testing URLs for Kubernetes specific scenarios should use an existing URL when there are URLs with no host defined in the env file with same port
odo devfile url command tests Testing URLs for OpenShift specific scenarios should create a automatically route on a openShift cluster
odo devfile url command tests Testing URLs for OpenShift specific scenarios should create a route on a openShift cluster without calling url create
odo devfile url command tests Testing URLs for OpenShift specific scenarios should create a url for a unsupported devfile component
odo devfile url command tests Testing URLs for OpenShift specific scenarios should error out when a host is provided with a route on a openShift cluster
odo devfile url command tests Testing URLs for OpenShift specific scenarios should list route and ingress urls with appropriate state

# watch command
odo devfile watch command tests when executing odo watch after odo push should listen for file changes
odo devfile watch command tests when executing odo watch after odo push with debug flag should be able to start a debug session after push with debug flag using odo watch and revert back after normal push
odo devfile watch command tests when executing odo watch after odo push with flag commands should listen for file changes
odo devfile watch command tests when executing odo watch after odo push with ignores flag should be able to ignore the specified file, .git and odo-file-index.json
odo devfile watch command tests when executing odo watch ensure that index information is updated by watch
odo devfile watch command tests when executing odo watch should listen for file changes with delay set to 0
odo devfile watch command tests when executing odo watch should show validation errors if the devfile is incorrect
odo devfile watch command tests when executing odo watch should use the index information from previous push operation
odo devfile watch command tests when executing watch without pushing a devfile component should error out on devfile flag
odo devfile watch command tests when executing watch without pushing a devfile component should fail
odo devfile watch command tests when running help for watch command should display the help
