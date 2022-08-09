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
	AbsPath string `yaml:"AbsPath"` // 存放的绝对路径
}

var LiveConfig Live

func LoadConfig() error {
	config := viper.New()
	config.AddConfigPath("./Api/setu")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	// 加载配置文件
	err := Config.PhraseYamlConfig("setu", "/Api/setu", defaultConfig, config)
	if err != nil {
		return err
	} // 直接返回错误信息到上一级
	// 加载live相关配置
	err = config.UnmarshalKey("live", &LiveConfig)
	if err != nil {
		fmt.Printf("setu/live Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
		return fmt.Errorf("setu/live Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
	}
	// 加载配置完成
	return nil
}
