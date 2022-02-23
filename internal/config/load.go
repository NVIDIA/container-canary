package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	canaryv1 "github.com/jacobtomlinson/containercanary/internal/apis/v1"
	"gopkg.in/yaml.v2"
)

func LoadValidatorFromURL(url string) (*canaryv1.Validator, error) {
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

func LoadValidatorFromFile(path string) (*canaryv1.Validator, error) {
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

func LoadValidatorFromBytes(b []byte) (*canaryv1.Validator, error) {
	var validator canaryv1.Validator

	err := yaml.Unmarshal(b, &validator)
	if err != nil {
		return nil, err
	}

	return &validator, nil
}
