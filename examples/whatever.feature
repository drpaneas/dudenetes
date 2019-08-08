Feature: The name of the feature

    Here write something to describe your feature
    with as many lines as you like

    Scenario: Initialize skuba structure for cluster deployment
        Given you have LB with IP "10.10.10.10"
        When you run "skuba cluster init --control-plane 10.10.10.10 my-cluster"
        Then a folder my-cluster should be created
            """
            ls -l my-cluster | grep my-cluster
            """
