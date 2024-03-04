package vspb

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General  *General   `yaml:"general"`
	Packages []*Package `yaml:"packages"`
}

type General struct {
	Path      string `yaml:"path"`
	RetryTime uint   `yaml:"retry_time"`
	Verbose   bool   `yaml:"verbose"`
}

type Package struct {
	Name string `yaml:"name"`

	// 插件名称
	// 会在构建的文件夹中搜索这些文件
	// 并复制到 path/plugins 路径中
	Provide string `yaml:"provide"`

	Address string `yaml:"address"`
	Version string `yaml:"version"`

	// 环境变量
	Env map[string]string `yaml:"env"`

	// 构建时执行的cmd
	// 使用该字段会忽略Tool以及Flag
	Run []string `yaml:"run"`

	// 不进行构建
	Skip bool `yaml:"skip"`
}

func ReadConfig(path string) (*Config, error) {
	if len(path) == 0 {
		path = "config.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	return config, err
}

func (c *Config) String() string {
	data, _ := yaml.Marshal(c)
	return string(data)
}

func (c *Config) Check() error {
	unique := make(map[string]struct{})
	for _, pkg := range c.Packages {
		if strings.Contains(pkg.Name, `/`) || strings.Contains(pkg.Name, `\`) {
			return fmt.Errorf("package name contain illegal characters: %s", pkg.Name)
		}

		_, ok := unique[pkg.Name]
		if ok {
			return fmt.Errorf("package name duplicate: %s", pkg.Name)
		}
		unique[pkg.Name] = struct{}{}
	}

	return nil
}
