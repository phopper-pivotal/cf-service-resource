package out

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Manifest struct {
	Data map[interface{}]interface{}
}

func NewManifest(manifestPath string) (Manifest, error) {
	yamlData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = yaml.Unmarshal(yamlData, &manifest.Data)
	if err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}
