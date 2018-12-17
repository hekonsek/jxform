package forms

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/jxform/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FormDefinition struct {
	Name         string        `yaml:"name"`
	Domain       string        `yaml:"domain"`
	Git          Git           `yaml:"git"`
	Environments []Environment `yaml:"environments"`
}

func (def FormDefinition) ResolveName() string {
	name, found := os.LookupEnv("NAME")
	if found {
		return name
	}
	return def.Name
}

type Git struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func (git Git) ResolveUsername() string {
	username, found := os.LookupEnv("GIT_USERNAME")
	if found {
		return username
	}
	return git.Username
}

func (git Git) ResolveToken() string {
	token, found := os.LookupEnv("GIT_TOKEN")
	if found {
		return token
	}
	return git.Token
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
	definitionName := definition.ResolveName()

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
		if strings.HasPrefix(line, definitionName+"\t") {
			hasCluster = true
			break
		}
	}
	if !hasCluster {
		eksCreateCommand := []string{"create", "cluster", "eks", "--cluster-name=" + definitionName, "--skip-installation=true"}
		eksCreateCommand = append(eksCreateCommand, fmt.Sprintf("--verbose=%t", verbose))
		if verbose {
			fmt.Printf("About to execute command: %s\n", append([]string{"jx"}, eksCreateCommand...))
		}
		err = util.NewExecs().Sout("jx", eksCreateCommand...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("EKS cluster %s already exists.\n", color.GreenString(definitionName))
	}

	helmList, err := util.NewExecs().Run("helm", "list")
	if err != nil {
		log.Fatal(err)
	}
	isJxInstalled := false
	for i, line := range helmList {
		if i == 0 {
			continue
		}
		if strings.HasPrefix(line, "jenkins-x\t") {
			isJxInstalled = true
			break
		}
	}
	if !isJxInstalled {
		installCommand := []string{"install", "--provider=eks", "-b",
			"--git-api-token=" + definition.Git.ResolveToken(), "--git-provider-url=" + definition.Git.Server, "--git-private=true", "--git-username=" + definition.Git.ResolveUsername(),
			"--default-environment-prefix=" + definitionName, "--no-default-environments=true", "--domain=" + definition.Domain}
		installCommand = append(installCommand, fmt.Sprintf("--verbose=%t", verbose))
		if verbose {
			fmt.Printf("About to execute command: %s\n", append([]string{"jx"}, installCommand...))
		}
		err = util.NewExecs().Sout("jx", installCommand...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Jenkins X already installed.\n")
	}

	return nil
}
