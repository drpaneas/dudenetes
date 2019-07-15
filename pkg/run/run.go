package run

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func containsArg(slice []string, arg string) (bool, int) {
	for index := range slice {
		if strings.Contains(slice[index], arg) {
			return true, index
		}
	}
	return false, 0
}

func SplitCmdInPipes(cmd string) (string, error) {
	var outputLeft string
	var output string
	var err error
	// It takes as granted there are pipes in the string
	slice := strings.Split(cmd, "|")

	for key, value := range slice {
		if key == 0 {
			outputLeft, err = Cmd(value)
			if err != nil {
				return outputLeft, err
			}
			continue
		}

		outputRight, err := Pipe(outputLeft, value)
		if err != nil {
			return outputRight, err
		}
		outputLeft = fmt.Sprint(outputRight)

		if key == (len(slice) - 1) {
			output = fmt.Sprint(outputLeft)
		}
	}
	return output, nil
}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

// Pipe blahblah
func Pipe(output string, pipe string) (string, error) {
	// Write the output we would like to pipe into a file
	err := writeToFile("tmp", output)
	if err != nil {
		log.Fatal(err)
	}

	c1 := exec.Command("cat", "tmp")

	slice := splitCmd(pipe)

	c2 := exec.Command(slice[0], slice[1:]...)

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
		err := fmt.Errorf("%s didn't return any result", pipe)
		deleteFile("tmp")
		return str, err
	}
	deleteFile("tmp")
	return str, nil

}

// SlowCmd executes a command and waits until it timeouts
func SlowCmd(cmd string, timeout int) (string, error) {
	slice := splitCmd(cmd)
	execute := exec.Command(slice[0], slice[1:]...)

	// Use a bytes.Buffer to get the stdout
	var stdout bytes.Buffer
	execute.Stdout = &stdout

	// Use a bytes.Buffer to ger the stderr
	var stderr bytes.Buffer
	execute.Stderr = &stderr

	execute.Start()

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- execute.Wait() }()

	// Start a timer
	timer := time.After(time.Duration(timeout) * time.Second)

	// Execute based on which channel we get the first message
	select {
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		execute.Process.Kill()
		return fmt.Sprintf("Command timed out (took more than %d seconds)", timer), fmt.Errorf("Needs more time than %d seconds", timer)
	case err := <-done:
		// Command completed before timeout. Print output and error if it exists.
		if err != nil {
			return string(fmt.Sprint(err) + ": " + stderr.String()), err
		}
	}
	return string(stdout.String()), nil
}

func SlowCmdDir(cmd string, timeout int, directory string) (string, error) {
	slice := splitCmd(cmd)
	execute := exec.Command(slice[0], slice[1:]...)
	execute.Dir = directory

	// Use a bytes.Buffer to get the stdout
	var stdout bytes.Buffer
	execute.Stdout = &stdout

	// Use a bytes.Buffer to ger the stderr
	var stderr bytes.Buffer
	execute.Stderr = &stderr

	execute.Start()

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- execute.Wait() }()

	// Start a timer
	timer := time.After(time.Duration(timeout) * time.Second)

	// Execute based on which channel we get the first message
	select {
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		execute.Process.Kill()
		return fmt.Sprintf("Command timed out (took more than %d seconds)", timer), fmt.Errorf("Needs more time than %d seconds", timer)
	case err := <-done:
		// Command completed before timeout. Print output and error if it exists.
		if err != nil {
			return string(fmt.Sprint(err) + ": " + stderr.String()), err
		}
	}
	return string(stdout.String()), nil
}

func splitCmd(cmd string) []string {
	slice := strings.Split(cmd, " ")
	if slice[len(slice)-1] == "" {
		slice = slice[:len(slice)-1]
	}
	if slice[0] == "" {
		slice = slice[1:]
	}
	for key, value := range slice {
		if value == "" {
			slice = append(slice[:key], slice[key+1:]...)
		}
	}
	return slice
}

// Cmd runs a command and has a default timeout for 2 seconds
func Cmd(cmd string) (string, error) {

	slice := splitCmd(cmd)
	execute := exec.Command(slice[0], slice[1:]...)

	// Use a bytes.Buffer to get the stdout
	var stdout bytes.Buffer
	execute.Stdout = &stdout

	// Use a bytes.Buffer to ger the stderr
	var stderr bytes.Buffer
	execute.Stderr = &stderr

	execute.Start()

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- execute.Wait() }()

	// Start a timer (minimm time to wait before bail out is 2 seconds)
	timer := time.After(5 * time.Second)

	// Execute based on which channel we get the first message
	select {
	case <-timer:
		// Timeout happened first, kill the process and print a message.
		execute.Process.Kill()
		return fmt.Sprintf("Command timed out (took more than %d seconds)", timer), fmt.Errorf("Needs more time than %d seconds", timer)
	case err := <-done:
		// Command completed before timeout. Print output and error if it exists.
		if err != nil {
			return string(fmt.Sprint(err) + ": " + stderr.String()), err
		}
	}
	return string(stdout.String()), nil
}

// CmdRetry will keep trying to run your Cmd() every 2 seconds
// until it either gets an expected result or timeout expires
func CmdRetry(cmd string, timeout int) (string, error) {

	retries := timeout / 2
	try := 0
	var err error
	var output string

	for {
		if strings.Contains(cmd, "|") {
			output, err = SplitCmdInPipes(cmd)
		} else {
			output, err = Cmd(cmd)
		}
		try++

		if err == nil {
			return output, nil
		}
		if try == retries {
			break
		}

	}
	return fmt.Sprintf("Maximum retries (%d) limit got reached)", retries), fmt.Errorf("Needs more time than %d seconds", timeout)
}

// LogError returns a pretty log
func LogError(cmd string, output string, err error) error {
	fmtOut := strings.Replace(output, "\n", "\n\t\t", -1)
	fmtErr := strings.Replace(fmt.Sprintf("%s", err), "\n", "\n\t\t", -1)
	return fmt.Errorf("\n\tFailed:\n\t-------\n\t%s\n\n\tOutput:\n\t-------\n\t%s\n\n\tError:\n\t-------\n\t%s", cmd, fmtOut, fmtErr)
}

func debugRun(cmd string) {
	slice := splitCmd(cmd)
	str := ""
	for index, element := range slice {
		fmt.Printf("%4d : %v\n", index, element)
		str = str + fmt.Sprintf("%s ", element)
	}
	fmt.Printf("The command is: %s\n", str)
}
