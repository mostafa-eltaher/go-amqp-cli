package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YamlConfig is exported.
type YamlConfig struct {
	Spec struct {
		Connection struct {
			Authentication struct {
				Type     string `yaml:"type"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"authentication"`
		} `yaml:"connection"`
		Sessions []struct {
			Links []struct {
				Role          string `yaml:"role"`
				Source        string `yaml:"source"`
				Target        string `yaml:"target"`
				InitialCredit int    `yaml:"initialCredit"`
			} `yaml:"links"`
		} `yaml:"sessions"`
	} `yaml:"spec"`
}

func main() {
	fmt.Println("Parsing YAML file")

	yamlConfig, err := parseConf()

	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	fmt.Printf("Result: %v\n", yamlConfig)
}

func parseConf() (*YamlConfig, error) {
	var fileName string
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")
	flag.Parse()

	if fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")

	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, errors.New("Error reading YAML file")
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)

	return &yamlConfig, err
}
