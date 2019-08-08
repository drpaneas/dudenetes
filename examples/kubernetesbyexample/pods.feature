Feature: Pods

    A pod is a collection of containers sharing a network
    and mount namespace and is the basic deployment in Kubernetes.
    All containers in a pod are scheduled on the same node.

    Scenario: Launch a pod
        Given you want to use the image "mhausenblas/simpleservice:0.5.0" and expose a HTTP API on port 9876
        When you execute "kubectl run sise --image=mhausenblas/simpleservice:0.5.0 --port=9876 --labels run=test"
        Then the pod with the label "run=test" should start and be ready for use within 60 seconds

    Scenario: Access the pod from within the cluster
        Given you have non-interactive access with the master node via "ssh sles@10.84.153.196"
        When you try to access its HTTP API requesting "/info" against its Pod IP address
        Then you should get a successful reply