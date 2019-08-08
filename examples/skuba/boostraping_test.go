package main

import (
	"os"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/drpaneas/dudenetes/pkg/run"
	"github.com/drpaneas/dudenetes/pkg/skuba"
)

var folder string

func youHaveDeployedTheRequiredInfrastructureForSUSECaaSP() error {

	// Load the current Terraform output of the working cluster
	err := skuba.LoadTF()
	if err != nil {
		return err
	}

	// Verify you have what you need
	err = skuba.Need(1, 1, 3) // LB, Masters, Workers
	if err != nil {
		return err
	}
	return nil
}

func youDo(arg1 string) error {
	arg1 = skuba.ReplaceVarsWithEnvs(arg1)
	output, err := run.Cmd(arg1)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func dirShouldBeCreatedContainingTheIPOfTheLoadbalancer(arg1, arg2 string) error {
	folder = arg1
	arg2 = skuba.ReplaceVarsWithEnvs(arg2)
	output, err := run.Cmd(arg2)
	if err != nil {
		return run.LogError(arg2, output, err)
	}
	return nil
}

func youRunWithATimeoutOfSeconds(arg1 string, arg2 int) error {
	arg1 = skuba.ReplaceVarsWithEnvs(arg1)
	output, err := run.SlowCmdDir(arg1, arg2, folder)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func afterConfiguringYourNewKubeconfigIntoThis(arg1 string) error {
	arg1 = strings.ReplaceAll(arg1, "$HOME", os.Getenv("HOME"))
	output, err := run.SlowCmdDir(arg1, 5, folder)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func theMasterMustBeReadyWithinSecondsTimeout(arg1 int, arg2 string) error {
	output, err := run.CmdRetry(arg2, arg1)
	if err != nil {
		return run.LogError(arg2, output, err)
	}
	return nil
}

func youRunSkubaNodeJoinWithSecTimeout(arg1 string, arg2 int) error {
	arg1 = skuba.ReplaceVarsWithEnvs(arg1)
	output, err := run.SlowCmdDir(arg1, arg2, folder)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func theNodeShouldBeReadyWithinSec(arg1 string, arg2 int) error {
	output, err := run.CmdRetry(arg1, arg2)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^you have deployed the required infrastructure for SUSE CaaSP$`, youHaveDeployedTheRequiredInfrastructureForSUSECaaSP)
	s.Step(`^you do "([^"]*)"$`, youDo)
	s.Step(`^"([^"]*)" dir should be created containing the IP of the loadbalancer "([^"]*)"$`, dirShouldBeCreatedContainingTheIPOfTheLoadbalancer)
	s.Step(`^you run "([^"]*)" with a timeout of (\d+) seconds$`, youRunWithATimeoutOfSeconds)
	s.Step(`^after configuring your new kubeconfig into this "([^"]*)"$`, afterConfiguringYourNewKubeconfigIntoThis)
	s.Step(`^the master must be ready within (\d+) seconds timeout "([^"]*)"$`, theMasterMustBeReadyWithinSecondsTimeout)
	s.Step(`^you run skuba node join "([^"]*)" with (\d+) sec timeout$`, youRunSkubaNodeJoinWithSecTimeout)
	s.Step(`^the node should be ready "([^"]*)" within (\d+) sec$`, theNodeShouldBeReadyWithinSec)
}
