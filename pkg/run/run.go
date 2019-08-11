package run

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

// CmdWithPipes returns a string with the output of the cmd execution and an error with its return code.
// It takes as granted there are pipes "|" in the command.
//
// CmdWithPipes splits the command into sub-commands (separated by pipes "|")
// and executes them in order (from left to right). When the first sub-command get executed
// it passes its output and the the next sub-command to Pipe(). This continues
// until there are no more sub-commands to be executed.
func CmdWithPipes(cmd string) (string, error) {
	var output string
	var err error

	// NOTE: it takes as granted cmd is consisted of pipes
	if !strings.Contains(cmd, "|") {
		err = fmt.Errorf("the command: \"%s\" does NOT include any pipes '|'", cmd)
		output = "(there is no output)"
		return output, err
	}
	subCommands := strings.Split(cmd, "|")

	// Loop through all the subCommands and execute them in pairs
	// passing the output of the left one (left side of the pipe)
	// as the last parameter to the right one (right side of the pipe)
	for i, subCommand := range subCommands {
		// Execute the first sub-command
		if i == 0 {
			output, err = Cmd(subCommand) // Execute the command (on the left side of the pipe)
			if err != nil {
				return output, err
			}
			continue
		}

		// Execute the next-subcommand (on the right side of the pipe)
		// piping the output of the previous one (on the left side)
		output, err = Pipe(output, subCommand)
		if err != nil {
			return output, err
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

// Pipe executes cmd passing input into it
// it returns the output of the execution and its error code
func Pipe(input, cmd string) (string, error) {
	cmdSlice := splitCmd(cmd)
	command := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	// stdIn will be connected to the command's standard input when the command starts.
	stdIn, err := command.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdIn.Close()

	// stdOut will be connected to the command's standard output when the command starts
	stdOut, err := command.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdOut.Close()

	if err = command.Start(); err != nil {
		return "", err
	}

	// Connect input with stdIn
	stdIn.Write([]byte(input))
	stdIn.Close()

	// Read the output (including errors) of this connection
	stdBytes, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", err
	}

	//  Wait for the command to exit
	command.Wait()

	output := string(stdBytes)
	if output == "" {
		err := fmt.Errorf("%s didn't return any result", cmd)
		return output, err
	}
	return output, nil
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

func SlowCmdEnv(cmd string, timeout int, env []string) (string, error) {
	slice := splitCmd(cmd)
	execute := exec.Command(slice[0], slice[1:]...)

	// Use the current Env
	execute.Env = os.Environ()

	// Export env vars
	execute.Env = append(execute.Env, env[0:]...)

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

// splitCmd returns a []string with trimmered whitespace
func splitCmd(cmd string) []string {
	cmdSlice := strings.Split(cmd, " ")
	var cmdSliceNoWhitespace []string
	for _, value := range cmdSlice {
		if value == "" {
			continue
		}
		cmdSliceNoWhitespace = append(cmdSliceNoWhitespace, value)
	}
	return cmdSliceNoWhitespace
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

// Cmd runs a command and has a default timeout for 2 seconds
func CmdEnv(cmd string, env []string) (string, error) {

	slice := splitCmd(cmd)
	execute := exec.Command(slice[0], slice[1:]...)

	// Use the current Env
	execute.Env = os.Environ()

	// Export env vars
	execute.Env = append(execute.Env, env[0:]...)

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
			output, err = CmdWithPipes(cmd)
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
