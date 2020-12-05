
package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GridConfig struct {
	TileSize int `yaml:"tile_size"`
}

type ConfigFile struct {
	Port string `yaml:"port"`
	Name string `yaml:"name"`
	Grid GridConfig `yaml:"grid"`
	MediaPath string `yaml:"media_path"`
	ExportPath string `yaml:"export_path"`
}

var Config ConfigFile

func GetConfig() ConfigFile {
	config := ConfigFile{}

	path, _ := filepath.Abs("./config.yaml")

	file, err := os.Open(path)
	FailOnError(err)

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	data := string(b)

	err = yaml.Unmarshal([]byte(data), &config)
	FailOnError(err)

	return config
}

func SetupConfig() {
	Config = GetConfig()
}