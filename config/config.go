package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

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
	//binary, _ := os.Executable()
	//root, _ := filepath.EvalSymlinks(filepath.Dir(binary))
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "")
	cfgName := root + "/config.yaml"

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err = yaml.Unmarshal(bytes, &BlockchainDataConfig); err != nil {
		fmt.Println(err.Error())
	}
}
