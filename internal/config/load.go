package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jacobtomlinson/containercanary/internal/apis/config"
	"gopkg.in/yaml.v2"
)

func LoadValidatorFromURL(url string) (*config.Validator, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return LoadValidatorFromBytes(body)
}

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

	return LoadValidatorFromBytes(yamlFile)
}

func LoadValidatorFromBytes(b []byte) (*config.Validator, error) {
	var validator config.Validator

	err := yaml.Unmarshal(b, &validator)
	if err != nil {
		return nil, err
	}

	return &validator, nil
}
