package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/drpaneas/dudenetes/pkg/run"

	"github.com/DATA-DOG/godog"
)

var checkVirt, dockerHelloWorldOutput string
var env []string

func youWantToUseMinikubeInYourMachine(arg1 string) error {
	expectedRuntime := arg1
	actualRuntime := fmt.Sprintf("%s", runtime.GOOS)
	if !strings.Contains(actualRuntime, expectedRuntime) {
		return fmt.Errorf("This machine is running %s instead of linux", runtime.GOOS)
	}
	return nil
}

func youCheckIfVirtualizationIsSupportedByRunning(arg1 string) error {
	cmd := arg1
	var err error
	checkVirt, err = run.Cmd(cmd)
	if err != nil {
		return err
	}
	return nil
}

func youShouldGetANonEmptyOutputAsAReply() error {
	if checkVirt != "" {
		return nil
	}
	return fmt.Errorf("The output is not empty. See:\n%s", checkVirt)
}

func youNeedDockerToBePresentInYourMachine() error {
	cmd := "which docker"
	output, err := run.Cmd(cmd)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youRun(arg1 string) error {
	var err error
	cmd := fmt.Sprintf("docker run hello-world")
	dockerHelloWorldOutput, err = run.SlowCmd(cmd, 60)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, dockerHelloWorldOutput, err)
}

func theOutputMessageShouldSay(arg1 string) error {
	expectedString := arg1
	actualString := fmt.Sprintf("%s", dockerHelloWorldOutput)
	if strings.Contains(actualString, expectedString) {
		return nil
	}
	return fmt.Errorf("Unexpected output. See:\n%s", actualString)
}

func youHaveMinikubeInstalledInYourMachine() error {
	cmd := "which minikube"
	output, err := run.Cmd(cmd)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youCanExecuteSudoCommandsWithoutAPassword() error {
	cmd := "sudo echo"
	output, err := run.Cmd(cmd)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func kubectlWorksWithoutProblems(arg1 string) error {
	cmd := arg1
	output, err := run.Cmd(cmd)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youPrepareTheEnvironmentVariablesAndFolders() error {
	env = []string{
		"MINIKUBE_WANTUPDATENOTIFICATION=false",
		"MINIKUBE_WANTREPORTERRORPROMPT=false",
		fmt.Sprintf("MINIKUBE_HOME=%s", os.Getenv("HOME")),
		"CHANGE_MINIKUBE_NONE_USER=true",
		fmt.Sprintf("KUBECONFIG=%s/.kube/config", os.Getenv("HOME")),
	}

	cmd := fmt.Sprintf("mkdir -p %s/.kube %s/.minikube", os.Getenv("MINIKUBE_HOME"), os.Getenv("MINIKUBE_HOME"))
	output, err := run.CmdEnv(cmd, env)
	if err != nil {
		return run.LogError(cmd, output, err)
	}

	cmd = fmt.Sprintf("touch %s", os.Getenv("KUBECONFIG"))
	output, err = run.CmdEnv(cmd, env)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youRunTheKubernetesComponentsOnTheHostAndNotInAVMUsing(arg1 string) error {
	cmd := arg1
	output, err := run.SlowCmdEnv(cmd, 300, env)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func checkingTheMinikubStatusShouldFinishSuccessfully(arg1 string) error {
	cmd := arg1
	output, err := run.CmdEnv(cmd, env)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youHaveAMinikubeCluster(arg1 string) error {
	cmd := arg1
	output, err := run.CmdEnv(cmd, env)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func youStopItAndDeleteIt(arg1, arg2 string) error {
	cmd := arg1
	output, err := run.SlowCmdEnv(cmd, 30, env)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	cmd = arg2
	output, err = run.SlowCmdEnv(cmd, 30, env)
	if err == nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func thereMustBeNoMore(arg1 string) error {
	cmd := arg1
	output, err := run.CmdEnv(cmd, env)
	// Reverse logic (error is expected to happen)
	if err != nil {
		return nil
	}
	return run.LogError(cmd, output, err)
}

func StartMinikube(s *godog.Suite) {
	s.Step(`^you want to use minikube in your "([^"]*)" machine$`, youWantToUseMinikubeInYourMachine)
	s.Step(`^you check if virtualization is supported by running "([^"]*)"$`, youCheckIfVirtualizationIsSupportedByRunning)
	s.Step(`^you should get a non empty output as a reply$`, youShouldGetANonEmptyOutputAsAReply)

	s.Step(`^you need docker to be present in your machine$`, youNeedDockerToBePresentInYourMachine)
	s.Step(`^you run "([^"]*)"$`, youRun)
	s.Step(`^the output message should say "([^"]*)"$`, theOutputMessageShouldSay)

	s.Step(`^you have minikube installed in your machine$`, youHaveMinikubeInstalledInYourMachine)
	s.Step(`^you can execute sudo commands without a password$`, youCanExecuteSudoCommandsWithoutAPassword)
	s.Step(`^kubectl works without problems "([^"]*)"$`, kubectlWorksWithoutProblems)
	s.Step(`^you prepare the environment variables and folders$`, youPrepareTheEnvironmentVariablesAndFolders)
	s.Step(`^you run the kubernetes components on the host and not in a VM using "([^"]*)"$`, youRunTheKubernetesComponentsOnTheHostAndNotInAVMUsing)
	s.Step(`^checking the minikub status "([^"]*)" should finish successfully$`, checkingTheMinikubStatusShouldFinishSuccessfully)
}

func StopMinikube(s *godog.Suite) {
	s.Step(`^you have a minikube cluster "([^"]*)"$`, youHaveAMinikubeCluster)
	s.Step(`^you stop it "([^"]*)" and delete it "([^"]*)"$`, youStopItAndDeleteIt)
	s.Step(`^there must be no more "([^"]*)"$`, thereMustBeNoMore)
}
