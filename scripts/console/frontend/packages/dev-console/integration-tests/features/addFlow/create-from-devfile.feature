@add-flow
Feature: Create Application from Devfile
              As a user, I want to deploy an application from git repo with devfile


        Background:
            Given user is at developer perspective
              And user has created or selected namespace "aut-addflow-devfile"
              And user is at Add page


        @regression
        Scenario: Deploy git workload with devfile from topology page: A-04-TC01
            Given user is at the Topology page
             When user right clicks on topology empty graph
              And user selects "Import from Git" option from Add to Project context menu
              And user enters Git Repo URL as "https://github.com/nodeshift-starters/devfile-sample" in Import from Git form
              And user enters workload name as "node-bulletin-board-1"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "node-bulletin-board-1" in topology page


        @regression
        Scenario: Create the workload from dev file: A-04-TC02
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/devfile-samples/devfile-sample-java-springboot-basic"
              And user enters workload name as "devfile-sample-java-springboot-basic"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-sample-java-springboot-basic" in topology page


        @regression @to-do
        Scenario: Create the workload from dev file: A-04-TC03
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/nodeshift-starters/devfile-sample"
              And user enters workload name as "node-bulletin-board"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "node-bulletin-board" in topology page


        @regression @to-do
        Scenario: Create the workload from dev file: A-04-TC04
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/devfile-samples/devfile-sample-code-with-quarkus"
              And user enters workload name as "devfile-sample-code-with-quarkus"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-sample-code-with-quarkus" in topology page


        @regression @to-do
        Scenario: Create the workload from dev file: A-04-TC05
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/devfile-samples/devfile-sample-python-basic"
              And user enters workload name as "devfile-sample-python-basic"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-sample-python-basic" in topology page


        @regression @to-do
        Scenario: Create the workload from dev file: A-04-TC06
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/devfile-samples/devfile-stacks-nodejs-react"
              And user enters workload name as "devfile-stacks-nodejs-react"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-stacks-nodejs-react" in topology page

        @regression @to-do
        Scenario: Create the workload from dev file: A-04-TC07
            Given user is at Import from Git form
             When user enters Git Repo URL as "https://github.com/devfile-samples/devfile-sample-dotnet60-basic"
              And user enters workload name as "devfile-sample-dotnet60-basic"
              And user clicks Create button on Add page
             Then user will be redirected to Topology page
              And user is able to see workload "devfile-sample-dotnet60-basic" in topology page
          
        # Below scenario to be removed after the tests are are updated
        # @smoke
        # Scenario: Create the sample workload from dev file: A-04-TC03
        #     Given user is at Import from Git page
        #      When user selects Try sample link
        #       And user clicks Create button on Devfile page
        #      Then user will be redirected to Topology page
        #       And user is able to see workload "devfile-sample" in topology page


        @regression @to-do
        Scenario: Create the Devfiles workload from Developer Catalog: A-04-TC03
            Given user is at Developer Catalog page
             When user clicks on Devfiles type
              And user clicks on Basic Python card
              And user clicks on Create Application on the side drawer
              And user enters Application name as "devfile-sample-python-basic-git-app" in DevFile page
              And user enters Name as "devfile-sample-python-basic-git1"
              And user clicks on Create
             Then user is able to see workload "devfile-sample-python-basic-git1" in topology page
