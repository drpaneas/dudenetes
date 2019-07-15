Feature: Bootstraping

    Bootstrapping the cluster is the initial process of starting up the cluster
    and defining which of the nodes are masters and which workers.
    For maximum automation of this process SUSE CaaS Platform uses the skuba package.

    The tests assumes you have skuba already available in your machine
    and you have already deployed the requred infrastructure.

    Scenario: Initialize the cluster
        Given you want to initialize a cluster called "my-cluster" using "10.84.154.72" as control-plane
        When you do the skuba init for this control-plane
        Then a folder named "my-cluster" should be generated

    Scenario: Bootstrap the master node
        Given you want to bootstrap a master node with IP "10.84.153.196"
        When you run skuba node bootstrap for this master and wait for 180 seconds
        Then an "admin.conf" will be created

    Scenario: Add 2 workers to the cluster
        Given you want to add a worker node with IP "10.84.73.167" and another one with "10.84.73.197"
        And copy the "admin.conf" into your "/home/tux/.kube/config" directory
        When you run skuba node join for both of them and wait for 180 seconds
        Then you should see 3 nodes when running "kubectl get nodes | grep Ready | wc -l" within 180 seconds