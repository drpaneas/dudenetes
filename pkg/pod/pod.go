package pod

import (
	"fmt"
	"strings"
	"time"

	"github.com/drpaneas/dudenetes/pkg/run"
)

func GetName(namespace string, label string) (string, error) {
	cmd := "kubectl -n " + namespace + " get --no-headers=true pods -l " + label
	output, err := run.SlowCmd(cmd, 10)
	if err != nil {
		return "Error: Couldn't find a relative pod.", run.LogError(cmd, output, err)
	}
	ss := strings.FieldsFunc(output, func(r rune) bool {
		if r == ' ' {
			return true
		}
		return false
	})
	podName := ss[0]
	return podName, nil
}

func GetStatus(podName string) (string, error) {

	// Describe the pod and grep to "Ready:"
	cmd := fmt.Sprintf("kubectl describe pod %s", podName)
	podDescribe, err := run.Cmd(cmd)
	if err != nil {
		return "Error: Couldn't describe the pod", run.LogError(cmd, podDescribe, err)
	}
	podStatus, err := run.Pipe(podDescribe, "grep Ready:")
	if err != nil {
		return "Error: There is no 'Ready:' status", run.LogError(cmd, podStatus, err)
	}
	return podStatus, nil
}

func IsReady(podName string, timeout int) bool {
	retries := timeout / 2
	try := 0
	for {
		podStatus, _ := GetStatus(podName)
		ready := strings.Contains(podStatus, "True")
		try++

		if ready {
			return true
		}
		if try == retries {
			return false
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}
