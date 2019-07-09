package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/drpaneas/dudenetes/pkg/command"
)

var LB string

func youHaveLBWithIP(arg1 string) error {
	LB = arg1
	return nil
}

func youRun(arg1 string) error {
	output, err := command.Run(arg1)
	if err != nil {
		return command.LogError(arg1, output)
	}
	return nil
}

func aFolderMyclusterShouldBeCreated(arg1 *gherkin.DocString) error {
	cmd := arg1.Content
	output, err := command.Run(cmd)
	if err != nil {
		return command.LogError(cmd, output)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^you have LB with IP "([^"]*)"$`, youHaveLBWithIP)
	s.Step(`^you run "([^"]*)"$`, youRun)
	s.Step(`^a folder my-cluster should be created$`, aFolderMyclusterShouldBeCreated)
}
