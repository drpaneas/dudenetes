package main

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/drpaneas/dudenetes/pkg/pod"
	"github.com/drpaneas/dudenetes/pkg/run"
)

var siseImage, siseName, podName, runMaster string
var sisePortAPI int

func youWantToUseTheImageAndExposeAHTTPAPIOnPort(arg1 string, arg2 int) error {
	siseImage = arg1
	sisePortAPI = arg2
	return nil
}

func youExecute(arg1 string) error {
	output, err := run.Cmd(arg1)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	return nil
}

func thePodWithTheLabelShouldStartAndBeReadyForUseWithinSeconds(arg1 string, arg2 int) error {

	podName, err := pod.GetName("default", arg1)
	if err != nil {
		return err
	}

	if pod.IsReady(podName, arg2) {
		return nil
	}

	return fmt.Errorf("Pod %s is not running", podName)

}

func youHaveNoninteractiveAccessWithTheMasterNodeVia(arg1 string) error {
	runMaster = arg1
	return nil
}

var nodeIP string

func youTryToAccessItsHTTPAPIRequestingAgainstItsPodIPAddress(arg1 string) error {
	cmd := "kubectl describe pod " + podName + " | grep IP"
	output, err := run.CmdWithPipes(cmd)
	if err != nil {
		return run.LogError(arg1, output, err)
	}
	// Extract the IP of the node that runs this container
	parts := strings.Split(output, ":")
	nodeIP = fmt.Sprintf(strings.TrimSpace(parts[1]))
	return nil
}

func youShouldGetASuccessfulReply() error {
	cmd := fmt.Sprintf("%s curl -s %s:%d/info", runMaster, nodeIP, sisePortAPI)
	output, err := run.CmdRetry(cmd, 30)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func PodFeatures(s *godog.Suite) {
	// Scenario: Launch a pod
	s.Step(`^you want to use the image "([^"]*)" and expose a HTTP API on port (\d+)$`, youWantToUseTheImageAndExposeAHTTPAPIOnPort)
	s.Step(`^you execute "([^"]*)"$`, youExecute)
	s.Step(`^the pod with the label "([^"]*)" should start and be ready for use within (\d+) seconds$`, thePodWithTheLabelShouldStartAndBeReadyForUseWithinSeconds)

	//  Scenario: Access the pod from within the cluster
	s.Step(`^you have non-interactive access with the master node via "([^"]*)"$`, youHaveNoninteractiveAccessWithTheMasterNodeVia)
	s.Step(`^you try to access its HTTP API requesting "([^"]*)" against its Pod IP address$`, youTryToAccessItsHTTPAPIRequestingAgainstItsPodIPAddress)
	s.Step(`^you should get a successful reply$`, youShouldGetASuccessfulReply)

}
