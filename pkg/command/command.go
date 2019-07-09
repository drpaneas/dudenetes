package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func containsArg(slice []string, arg string) (bool, int) {
	for index := range slice {
		if strings.Contains(slice[index], arg) {
			return true, index
		}
	}
	return false, 0
}

// Run executes a command amd returns the output and the return code
func Run(cmd string) (string, error) {

	slice := strings.Split(cmd, " ")
	containsPipe, index := containsArg(slice, "|")

	// Works only for 1 pipe at the moment
	if containsPipe {

		// split slice into two parts. One before the pipe and one after
		beforePipe := slice[0 : index-1]
		afterPipe := slice[index+1:]

		// See https://golang.org/pkg/os/exec/#Cmd.StdinPipe
		c1 := exec.Command(beforePipe[0], beforePipe[1:]...)
		c2 := exec.Command(afterPipe[0], afterPipe[1:]...)
		r, w := io.Pipe()
		c1.Stdout = w
		c2.Stdin = r
		var b2 bytes.Buffer
		c2.Stdout = &b2
		c1.Start()
		c2.Start()
		c1.Wait()
		w.Close()
		c2.Wait()
		str := ""
		str = b2.String()
		if str == "" {
			err := fmt.Errorf("%s didn't return any result", afterPipe[1])
			return str, err
		}
		return str, nil
	}

	execute := exec.Command(slice[0], slice[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	execute.Stdout = &out
	execute.Stderr = &stderr
	err := execute.Run()
	if err != nil {
		return string(fmt.Sprint(err) + ": " + stderr.String()), err
	}
	return string(out.String()), nil
}

// LogError returns a pretty log
func LogError(cmd string, output string) error {
	fmtOut := strings.Replace(output, "\n", "\n\t\t", -1)
	return fmt.Errorf("\tFailed:\n\t-------\n\t%s\n\n\tOutput:\n\t-------\n\t%s", cmd, fmtOut)
}

func debugRun(cmd string) {
	slice := strings.Split(cmd, " ")
	str := ""
	for index, element := range slice {
		fmt.Printf("%4d : %v\n", index, element)
		str = str + fmt.Sprintf("%s ", element)
	}
	fmt.Printf("The command is: %s\n", str)
}
