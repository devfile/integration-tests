@odo commands


d odo devfile catalog command tests When checking catalog for installed services should succeed


odo devfile catalog command tests When executing catalog describe component with -o json should display a valid JSON


odo devfile catalog command tests When executing catalog describe component with a component name with multiple components should print multiple devfiles from different registries


d odo devfile catalog command tests When executing catalog list components should list all supported devfile components

odo devfile catalog command tests When executing catalog list components should list components successfully even with an invalid kubeconfig path or path points to existing directory


Feature: Catalog command list supported devfile components
              As a user, I want to list all supported devfile components

        Background:
            Given user logged in as developer

        @catalog
        Scenario: odo devfile catalog command tests
             When executing catalog list components
             Then should list all supported devfile components

        @catalog
        Scenario: Create the workload from dev file: A-04-TC02
            Given user is at Import from Devfile page
             When user enters Git Repo url "https://github.com/redhat-developer/devfile-sample" in Devfile Page
              And user enters Name as "node-bulletin-board" in DevFile page
              And user clicks Create button on Devfile page
             Then user will be redirected to Topology page
              And user is able to see workload "node-bulletin-board" in topology page


        @regression
        Scenario: Create the sample workload from dev file: A-04-TC03
            Given user is at Import from Devfile page
             When user selects Try sample link
              And user clicks Create button on Devfile page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-sample" in topology page
