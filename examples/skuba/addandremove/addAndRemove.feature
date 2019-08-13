Feature: Bootstraping

    Bootstrapping the cluster is the initial process of starting up the cluster
    and defining which of the nodes are masters and which workers.
    For maximum automation of this process SUSE CaaS Platform uses the skuba package.

    The tests assumes you have skuba already available in your machine
    and you have already deployed the requred infrastructure along with
    the SSH-agent running from the terminal you are issuing dudenetes commands.

    Scenario: Initialize the cluster
        Given there a deployed infrastructure for 1 lb, 3 masters and 3 workers
        When you do "skuba cluster init --control-plane $loadbalancer my-cluster"
        Then "my-cluster" dir should be created containing the IP of the loadbalancer "grep -r $loadbalancer my-cluster"

    Scenario: Bootstrap the master node
        Given you run "skuba -v 5 node bootstrap --user sles --sudo --target $master1 master-1" with a timeout of 500 seconds
        And after configuring your new kubeconfig into this "cp admin.conf $HOME/.kube/config"
        Then the master must be ready within 500 seconds timeout "kubectl get nodes |  grep master-1 | grep --invert-match NotReady | grep Ready"

    Scenario: Join 1 worker
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker1 worker-1" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-1 | grep --invert-match NotReady | grep Ready" within 500 sec

    Scenario: Add another master node, remove it, and then add another master node
        Given you run "skuba -v 5 node join --role master --user sles --sudo --target $master2 master-2" with a timeout of 500 seconds
        Then the master must be ready within 500 seconds timeout "kubectl get nodes |  grep master-2 | grep --invert-match NotReady | grep Ready"
        And now you must have two ready masters "kubectl get nodes |  grep master | grep --invert-match NotReady | grep Ready | wc -l | grep 2"
        When you remove this master node "skuba node remove master-2 --drain-timeout 5s"
        Then there must be only one master at your cluster "kubectl get nodes | grep master | grep --invert-match NotReady | grep Ready | wc -l | grep 1"
        Given you run "skuba -v 5 node join --role master --user sles --sudo --target $master3 master-3" with a timeout of 500 seconds
        Then the master must be ready within 500 seconds timeout "kubectl get nodes |  grep master-3 | grep --invert-match NotReady | grep Ready"
        And now you must have two ready masters "kubectl get nodes |  grep master | grep --invert-match NotReady | grep Ready | wc -l | grep 2"

    Scenario: Join another worker node, remove it, and then add another worker node
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker2 worker-2" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-2 | grep --invert-match NotReady | grep Ready" within 500 sec
        And now you must have two ready workers "kubectl get nodes |  grep worker | grep --invert-match NotReady | grep Ready | wc -l | grep 2"
        When you remove this worker node "skuba node remove worker-2 --drain-timeout 5s"
        Then there must be only one worker at your cluster "kubectl get nodes | grep worker | grep --invert-match NotReady | grep Ready | wc -l | grep 1"
        When you run skuba node join "skuba -v 5 node join --role worker --user sles --sudo --target $worker3 worker-3" with 500 sec timeout
        Then the node should be ready "kubectl get nodes | grep worker-3 | grep --invert-match NotReady | grep Ready" within 500 sec
        And now you must have two ready workers "kubectl get nodes |  grep worker | grep --invert-match NotReady | grep Ready | wc -l | grep 2"



