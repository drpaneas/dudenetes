Feature: Installing Minikube

    https://kubernetes.io/docs/tasks/tools/install-minikube/

    Scenario: Initialize the cluster
        Given you want to use minikube in your "linux" machine
        When you check if virtualization is supported by running "grep vmx /proc/cpuinfo"
        Then you should get a non empty output as a reply

    Scenario: Check docker status
        Given you need docker to be present in your machine
        When you run "docker run hello-world"
        Then the output message should say "Hello from Docker"

    Scenario: Start minikube with 'none' driver
        Given you have minikube installed in your machine
        And you can execute sudo commands without a password
        And kubectl works without problems "kubectl version --client"
        When you prepare the environment variables and folders
        And you run the kubernetes components on the host and not in a VM using "sudo -E minikube start --vm-driver=none"
        Then checking the minikub status "minikube status" should finish successfully

    Scenario: Stop and Delete minikube deployment
        Given you have a minikube cluster "minikube status"
        When you stop it "minikube stop" and delete it "minikube delete"
        Then there must be no more "minikube status"