package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/DATA-DOG/godog"
)

var loadBalancerIP, clusterLocation string

// Feature 1

func youHaveALoadbalancerUpAndRunningWith(arg1 string) error {
	loadBalancerIP = arg1
	if loadBalancerIP == "" {
		log.Fatal("There is no IP for Load Balancer")
	}
	return nil
}

func youInitializeASkubaStructureForDeployment(arg1 string) error {
	clusterLocation = arg1
	cmd := exec.Command("skuba", "cluster", "init", "--control-plane", loadBalancerIP, clusterLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func aFolderShouldBeCreated(arg1 string) error {
	cmd := exec.Command("ls", "-l", clusterLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil

}

// Feature 2
var files []os.FileInfo

func thereIsFolderCalled(arg1 string) error {
	cmd := exec.Command("ls", "-l", arg1)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func youBrowseInsideOfIt() error {
	var err error
	files, err = ioutil.ReadDir("./my-cluster")
	if err != nil {
		return err
	}
	return nil
}

func itShouldHaveListings(arg1 int) error {
	numberOfExpectedFiles := arg1
	numberOfActualFiles := len(files)
	if numberOfExpectedFiles != numberOfActualFiles {
		errMsg := fmt.Sprintf("It found %d files instead of %d", numberOfActualFiles, numberOfExpectedFiles)
		return errors.New(errMsg)
	}
	return nil
}

func theseWillBeDirectoriesAnd(arg1 int, arg2, arg3 string) error {
	numOfDirs := 0
	for _, value := range files {
		if value.Name() == arg2 || value.Name() == arg3 {
			numOfDirs++
		}
	}
	if numOfDirs != arg1 {
		errMsg := fmt.Sprintf("It found %d directories instead of %d", numOfDirs, arg1)
		return errors.New(errMsg)
	}
	return nil
}

func fileCalled(arg1 int, arg2 string) error {
	numOfFiles := 0
	for _, value := range files {
		if value.Name() == arg2 {
			numOfFiles++
		}
	}
	if numOfFiles != arg1 {
		errMsg := fmt.Sprintf("It found %d files instead of %d", numOfFiles, arg1)
		return errors.New(errMsg)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^you have a load-balancer up and running with "([^"]*)"$`, youHaveALoadbalancerUpAndRunningWith)
	s.Step(`^you initialize a skuba structure for "([^"]*)" deployment$`, youInitializeASkubaStructureForDeployment)
	s.Step(`^a folder "([^"]*)" should be created$`, aFolderShouldBeCreated)

	s.Step(`^there is folder called "([^"]*)"$`, thereIsFolderCalled)
	s.Step(`^you browse inside of it$`, youBrowseInsideOfIt)
	s.Step(`^it should have (\d+) listings$`, itShouldHaveListings)
	s.Step(`^these will be (\d+) directories "([^"]*)" and "([^"]*)"$`, theseWillBeDirectoriesAnd)
	s.Step(`^(\d+) file called "([^"]*)"$`, fileCalled)

}
