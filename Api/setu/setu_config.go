package setu

import (
	"TsubakiApi/Config"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
)

//go:embed default_config.yaml
var defaultConfig string

//Live 三次元涩图路径
type Live struct {
	absPath      string `yaml:"absPath"`
	relativePath string `yaml:"relativePath"`
}

var LiveConfig Live

func LoadConfig() error {
	config := viper.New()
	config.AddConfigPath("./Api/setu")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	//加载配置文件
	err := Config.PhraseYamlConfig("setu", "/Api/setu", defaultConfig, config)
	if err != nil {
		return err
	} // 直接返回错误信息到上一级
	err = config.UnmarshalKey("live", &LiveConfig)
	if err != nil {
		fmt.Printf("Api/setu/live Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
		return fmt.Errorf("Api/setu/live Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
	}
	return nil
}
