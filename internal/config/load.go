package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jacobtomlinson/containercanary/internal/apis/config"
	"gopkg.in/yaml.v2"
)

func LoadValidatorFromFile(path string) (*config.Validator, error) {
	filename, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("no such file %s", filename)
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var validator config.Validator

	err = yaml.Unmarshal(yamlFile, &validator)
	if err != nil {
		return nil, err
	}
	return &validator, nil
}
