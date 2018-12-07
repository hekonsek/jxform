package forms

import (
	"fmt"
	"github.com/hekonsek/jxform/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

	eksCreateCommand := []string{"create", "cluster", "eks", "--skip-installation=true"}
	eksCreateCommand = append(eksCreateCommand, fmt.Sprintf("--verbose=%t", verbose))
	err = util.NewExecs().Sout("jx", eksCreateCommand...)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
