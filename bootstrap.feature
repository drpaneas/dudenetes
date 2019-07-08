Feature: Bootstrap 3 node kubernetes cluster
    In order to bootstrap a 3 node kubernetes cluster
    As a user having the infrastructure ready with Terraform
    I need run skuba commands against this infrastructure

Scenario: Initialize skuba structure for cluster deployment
    Given you have a load-balancer up and running with "10.10.10.10"
    When you initialize a skuba structure for "my-cluster" deployment
    Then a folder "my-cluster" should be created

Scenario: Verify the skuba structure
    Given there is folder called "my-cluster"
    When you browse inside of it
    Then it should have 3 listings
    And these will be 2 directories "addons" and "kubeadm-join.conf.d"
    And 1 file called "kubeadm-init.conf"