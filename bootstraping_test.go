package main

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/drpaneas/dudenetes/pkg/run"
	"github.com/drpaneas/dudenetes/pkg/skuba"
)

var cluster, lb, master, worker1, worker2 string

func youWantToInitializeAClusterCalledUsingAsControlplane(arg1, arg2 string) error {
	cluster = arg1
	lb = arg2
	return nil
}

func youDoTheSkubaInitForThisControlplane() error {
	err := skuba.InitControlPlane(lb, cluster)
	if err != nil {
		return err
	}
	return nil
}

func aFolderNamedShouldBeGenerated(arg1 string) error {
	cmd := fmt.Sprintf("ls -l %s | grep -r %s", arg1, lb)
	output, err := run.SplitCmdInPipes(cmd)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func youWantToBootstrapAMasterNodeWithIP(arg1 string) error {
	master = arg1
	return nil
}

func youRunSkubaNodeBootstrapForThisMasterAndWaitForSeconds(arg1 int) error {
	err := skuba.NodeBootstrap(master, cluster, "master-1", arg1)
	if err != nil {
		return err
	}
	return nil
}

func anWillBeCreated(arg1 string) error {
	cmd := fmt.Sprintf("ls -l %s | grep %s", cluster, arg1)
	output, err := run.SplitCmdInPipes(cmd)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func youWantToAddAWorkerNodeWithIPAndAnotherOneWith(arg1, arg2 string) error {
	worker1 = arg1
	worker2 = arg2
	return nil
}

func copyTheIntoYourDirectory(arg1, arg2 string) error {
	cluster = "my-cluster"
	cmd := fmt.Sprintf("cp %s/%s %s", cluster, arg1, arg2)
	output, err := run.Cmd(cmd)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func youRunSkubaNodeJoinForBothOfThemAndWaitForSeconds(arg1 int) error {
	cluster = "my-cluster"
	// Join worker-1
	err := skuba.NodeJoinWorker(worker1, cluster, "worker-1", arg1)
	if err != nil {
		return err
	}
	// Join worker-2
	err = skuba.NodeJoinWorker(worker2, cluster, "worker-2", arg1)
	if err != nil {
		return err
	}
	return nil
}

func youShouldSeeNodesWhenRunningWithinSeconds(arg1 int, arg2 string, arg3 int) error {
	cmd := arg2
	expectedReadyNodes := fmt.Sprintf("%d", arg1)
	output, err := run.CmdRetry(cmd, arg3)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	if !strings.Contains(output, expectedReadyNodes) {
		return fmt.Errorf("Only %s workers out of %s are Ready", output, expectedReadyNodes)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^you want to initialize a cluster called "([^"]*)" using "([^"]*)" as control-plane$`, youWantToInitializeAClusterCalledUsingAsControlplane)
	s.Step(`^you do the skuba init for this control-plane$`, youDoTheSkubaInitForThisControlplane)
	s.Step(`^a folder named "([^"]*)" should be generated$`, aFolderNamedShouldBeGenerated)

	s.Step(`^you want to bootstrap a master node with IP "([^"]*)"$`, youWantToBootstrapAMasterNodeWithIP)
	s.Step(`^you run skuba node bootstrap for this master and wait for (\d+) seconds$`, youRunSkubaNodeBootstrapForThisMasterAndWaitForSeconds)
	s.Step(`^an "([^"]*)" will be created$`, anWillBeCreated)

	s.Step(`^you want to add a worker node with IP "([^"]*)" and another one with "([^"]*)"$`, youWantToAddAWorkerNodeWithIPAndAnotherOneWith)
	s.Step(`^copy the "([^"]*)" into your "([^"]*)" directory$`, copyTheIntoYourDirectory)
	s.Step(`^you run skuba node join for both of them and wait for (\d+) seconds$`, youRunSkubaNodeJoinForBothOfThemAndWaitForSeconds)
	s.Step(`^you should see (\d+) nodes when running "([^"]*)" within (\d+) seconds$`, youShouldSeeNodesWhenRunningWithinSeconds)

}
