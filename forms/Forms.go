package forms

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/jxform/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type FormDefinition struct {
	Name         string        `yaml:"name"`
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name string `yaml:"name"`
}

func Provision(verbose bool) error {
	definitionFile, err := ioutil.ReadFile("example.yml")
	if err != nil {
		return err
	}

	definition := FormDefinition{}
	err = yaml.Unmarshal(definitionFile, &definition)
	if err != nil {
		return err
	}

	eksGetCommand := []string{"get", "eks"}
	if verbose {
		fmt.Printf("About to execute command: %s\n", append([]string{"jx"}, eksGetCommand...))
	}
	out, err := util.NewExecs().Run("jx", eksGetCommand...)
	if err != nil {
		log.Fatal(err)
	}
	hasCluster := false
	for _, line := range out {
		if strings.HasPrefix(line, definition.Name+"\t") {
			hasCluster = true
			break
		}
	}
	if !hasCluster {
		eksCreateCommand := []string{"create", "cluster", "eks", "--cluster-name=" + definition.Name, "--skip-installation=true"}
		eksCreateCommand = append(eksCreateCommand, fmt.Sprintf("--verbose=%t", verbose))
		if verbose {
			fmt.Printf("About to execute command: %s\n", append([]string{"jx"}, eksCreateCommand...))
		}
		err = util.NewExecs().Sout("jx", eksCreateCommand...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Cluster %s already exists.\n", color.GreenString(definition.Name))
	}

	return nil
}
