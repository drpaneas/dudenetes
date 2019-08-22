package skuba

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/drpaneas/dudenetes/pkg/run"
	"github.com/joho/godotenv"
)

func InitControlPlane(IP string, folder string) error {
	cmd := fmt.Sprintf("skuba cluster init --control-plane %s %s", IP, folder)
	output, err := run.Cmd(cmd)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func NodeBootstrap(IP string, folder string, masterNode string, timeout int) error {
	cmd := fmt.Sprintf("skuba node bootstrap --user sles --sudo --target %s %s", IP, masterNode)
	// Wait max 3 minutes to bootstrap the master node
	output, err := run.SlowCmdDir(cmd, timeout, folder)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func NodeJoinWorker(IP string, folder string, workerNode string, timeout int) error {
	cmd := fmt.Sprintf("skuba node join --role worker --user sles --sudo --target %s %s", IP, workerNode)
	// Wait max 3 minutes to bootstrap the master node
	output, err := run.SlowCmdDir(cmd, timeout, folder)
	if err != nil {
		return run.LogError(cmd, output, err)
	}
	return nil
}

func LoadTF() error {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		return errors.New("No .env file found")
	}
	return nil
}

func envExists(key string) bool {
	_, ok := os.LookupEnv(key)
	if ok {
		return true
	}
	return false
}

func Need(lb, masters, workers int) error {
	numMasters := 1
	for {
		envVar := fmt.Sprintf("master%d", numMasters)
		if envExists(envVar) {
			numMasters++
		} else {
			numMasters--
			break
		}
	}

	numWorkers := 1
	for {
		envVar := fmt.Sprintf("master%d", numWorkers)
		if envExists(envVar) {
			numWorkers++
		} else {
			numWorkers--
			break
		}
	}

	numLB := 0
	if envExists("loadbalancer") {
		numLB++
	}

	if (numLB >= lb) && (numMasters >= masters) && (numWorkers >= workers) {
		return nil
	}

	err := fmt.Errorf("\nRequired : LB %d, Masters %d, Workers %d\nAvailable: LB %d, Masters %d, Workers %d", lb, masters, workers, numLB, numMasters, numWorkers)
	return err

}

func FindEnvVar(str, before, after string) string {
	a := strings.SplitAfterN(str, before, 2)
	b := strings.SplitAfterN(a[len(a)-1], after, 2)
	if 1 == len(b) {
		return b[0]
	}
	return b[0][0 : len(b[0])-len(after)]
}

// ReplaceVarsWithEnvs works only with one variable
func ReplaceVarsWithEnvs(cmd string) string {
	if strings.Contains(cmd, "$") {
		v := FindEnvVar(cmd, "$", " ")
		cmd = strings.Replace(cmd, fmt.Sprintf("$%s", v), os.Getenv(v), -1)
	}
	return cmd
}
