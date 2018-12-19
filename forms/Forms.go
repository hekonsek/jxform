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
import "github.com/pkg/errors"

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
	Owner    string `yaml:"owner"`
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

func (git Git) ResolveOwner() string {
	owner, found := os.LookupEnv("GIT_OWNER")
	if found {
		return owner
	}
	return git.Owner
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
			return err
		}
	} else {
		fmt.Printf("EKS cluster %s already exists.\n", color.GreenString(definitionName))
	}

	helmList, err := util.NewExecs().Run("helm", "list")
	if err != nil {
		if !strings.Contains(helmList[0], "could not find tiller") && !strings.Contains(helmList[0], "configmaps is forbidden") {
			return errors.Wrap(err, strings.Join(helmList, "\n"))
		}
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

	existingEnvironments, err := util.NewExecs().Run("jx", "get", "env")
	if err != nil {
		return err
	}
	for _, environment := range definition.Environments {
		environmentExists := false
		for i, existingEnvironment := range existingEnvironments {
			if i == 0 {
				continue
			}
			if strings.HasPrefix(existingEnvironment, environment.Name+" ") {
				environmentExists = true
				break
			}
		}
		if environmentExists {
			fmt.Printf("Environment %s already exists.\n", color.GreenString(environment.Name))
		} else {
			createEnvCommand := []string{"create", "env", "-b", "--name=" + strings.ToLower(environment.Name), "--label=" + strings.ToLower(environment.Name),
				"--promotion=Auto",
				"--git-provider-url=" + definition.Git.Server, "--git-username=" + definition.Git.ResolveUsername(), "--git-owner=" + definition.Git.ResolveOwner(),
				"--git-private=true", "--git-api-token=" + definition.Git.ResolveToken()}
			createEnvCommand = append(createEnvCommand, fmt.Sprintf("--verbose=%t", verbose))
			if verbose {
				fmt.Printf("About to execute command: %s\n", append([]string{"jx"}, createEnvCommand...))
			}
			err = util.NewExecs().Sout("jx", createEnvCommand...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
