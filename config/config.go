package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var BlockchainDataConfig blockchainDataConfig

type blockchainDataConfig struct {
	ServerVersion string `yaml:"serverversion"`

	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
}

func init() {
	binary, _ := os.Executable()
	root, _ := filepath.EvalSymlinks(filepath.Dir(binary))
	cfgName := root + "/config/config.yaml"

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err.Error())
	}
	if err = yaml.Unmarshal(bytes, &BlockchainDataConfig); err != nil {
		panic(err.Error())
	}
}
