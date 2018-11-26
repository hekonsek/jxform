package forms

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FormDefinition struct {
	Name         string        `yaml:"name"`
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name string `yaml:"name"`
}

func Provision() error {
	definitionFile, err := ioutil.ReadFile("example.yml")
	if err != nil {
		return err
	}

	definition := FormDefinition{}
	err = yaml.Unmarshal(definitionFile, &definition)
	if err != nil {
		return err
	}

	fmt.Println(definition)

	return nil
}
