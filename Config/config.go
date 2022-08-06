// Package Config: 基于viper的配置读取
package Config

import (
	"TsubakiApi/Utils"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

//go:embed default_config.yaml
var defaultConfig string

//LogConfig 日志配置相关
type LogConfig struct {
	Filename           string `yaml:"Filename"`           //日志路径
	MaxSize            int    `yaml:"MaxSize"`            //日志分割大小 (MB)
	MaxBackups         int    `yaml:"MaxBackups"`         //保留旧日志文件的最大个数
	MaxAge             int    `yaml:"MaxAge"`             //保留旧日志文件的最长天数
	Compress           bool   `yaml:"Compress"`           //是否压缩旧日志文件
	DoNotDeleteLogFile bool   `yaml:"DoNotDeleteLogFile"` //为真时会保留所有日志文件，MaxBackup和MaxAge选项失效
}

type ServerConfig struct {
	Port int `yaml:"Port"` //端口号
}

var Log LogConfig       // 日志配置
var Server ServerConfig // 服务器配置

//Parse 解析配置文件
func Parse() {
	// 配置文件读取
	viper.SetConfigName("config")                // 配置文件名
	viper.SetConfigType("yaml")                  // 配置文件后缀
	viper.AddConfigPath("./Config")              // 配置文件路径
	if err := viper.ReadInConfig(); err != nil { //配置文件读取
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 没有发现配置文件
			fmt.Println("config.Parse:  没有发现配置文件，将自动生成一个新的配置文件。")
			path, err := os.Getwd() //获取执行路径
			if err != nil {
				fmt.Println("config.Parse:  获取当前执行路径出错。错误信息:", err)
				panic("config.Parse:  获取当前执行路径出错。")
			}
			path = path + "/Config" //配置文件路径
			if ok, _ := Utils.HasDir(path); !ok {
				//路径不存在
				Utils.CreateDir(path) //创建配置文件夹 ./Config
			}
			// 创建配置文件
			configFile, _ := os.OpenFile("./Config/config.yaml", os.O_CREATE|os.O_RDWR, 0777)
			_, _ = configFile.WriteString(defaultConfig) // 将默认的配置写入文件
			_ = configFile.Close()
			if err = viper.ReadInConfig(); err != nil {
				fmt.Println("config.Parse:  配置文件读取出错。错误信息:", err)
				panic("config.Parse:  配置文件读取出错。")
			}
		} else {
			fmt.Println("config.Parse:  配置文件读取出错。")
			panic("config.Parse:  配置文件读取出错。")
		}
	}
	//log 配置文件解析
	err := viper.UnmarshalKey("log", &Log) // 解析log配置
	if err != nil {
		fmt.Println("config.Parse:  配置文件log部分格式错误")
		panic("config.Parse:  配置文件log部分格式错误")
	}
	//Server 配置文件解析
	err = viper.UnmarshalKey("Server", &Server) //解析server配置
	if err != nil {
		fmt.Println("config.Parse:  配置文件log部分格式错误")
		panic("config.Parse:  配置文件log部分格式错误")
	}
}
