package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/jacobtomlinson/containercanary/internal/apis/config"
	"gopkg.in/yaml.v2"
)

func LoadValidatorFromFile(path string) (*config.Validator, error) {
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var validator config.Validator

	err = yaml.Unmarshal(yamlFile, &validator)
	if err != nil {
		return nil, err
	}
	return &validator, nil
}
