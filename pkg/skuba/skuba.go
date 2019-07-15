package skuba

import (
	"fmt"

	"github.com/drpaneas/dudenetes/pkg/run"
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
