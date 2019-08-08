Feature: Bootstraping

    Bootstrapping the cluster is the initial process of starting up the cluster
    and defining which of the nodes are masters and which workers.
    For maximum automation of this process SUSE CaaS Platform uses the skuba package.

    The tests assumes you have skuba already available in your machine
    and you have already deployed the requred infrastructure along with
    the SSH-agent running from the terminal you are issuing dudenetes commands.

    Scenario: Initialize the cluster
        Given you have deployed the required infrastructure for SUSE CaaSP
        When you do "skuba cluster init --control-plane $loadbalancer my-cluster"
        Then "my-cluster" dir should be created containing the IP of the loadbalancer "grep -r $loadbalancer my-cluster"

    Scenario: Bootstrap the master node
        Given you run "skuba -v 5 node bootstrap --user sles --sudo --target $master1 master-1" with a timeout of 500 seconds
        And after configuring your new kubeconfig into this "cp admin.conf $HOME/.kube/config"
        Then the master must be ready within 500 seconds timeout "kubectl get nodes |  grep master-1 | grep --invert-match NotReady | grep Ready"

    Scenario: Join the workers
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker1 worker-1" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-1 | grep --invert-match NotReady | grep Ready" within 180 sec
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker2 worker-2" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-2 | grep --invert-match NotReady | grep Ready" within 180 sec
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker3 worker-3" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-3 | grep --invert-match NotReady | grep Ready" within 180 sec

