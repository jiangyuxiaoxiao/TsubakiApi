package Config

import (
	"github.com/spf13/viper"
)

/*
Parse 解析配置文件
*/
func Parse() {
	viper.SetConfigName("config")                // 配置文件名
	viper.SetConfigType("yaml")                  // 配置文件后缀
	viper.AddConfigPath("./Config")              // 配置文件路径
	viper.AddConfigPath("../config")             //配置文件路径
	if err := viper.ReadInConfig(); err != nil { //配置文件读取
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
		}
	}

}
