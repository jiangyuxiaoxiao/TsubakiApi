// Package Config: 基于viper的配置读取
package Config

import (
	"TsubakiApi/Utils"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
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

//PhraseYamlConfig Viper解析yaml文件。适用于所有Api插件各自解析路径。
//当配置文件路径不存在时，将根据传入的值创建对应文件夹下的configFile。
//当文件夹都不存在时，也会新建文件夹。
//pack 包名 path 到插件的相对路径名,不加. defaultFile 默认配置文件string。
//path举例 执行路径名为 . 插件所在文件夹路径名./Api/setu 则path = /Api/setu
func PhraseYamlConfig(pack string, path string, defaultFile string, config *viper.Viper) error {
	if err := config.ReadInConfig(); err != nil { //配置文件读取
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 没有发现配置文件
			fmt.Printf("%s.ParseConfig:  没有发现[%s]配置文件，将自动生成一个新的配置文件。\n", pack, pack)
			ProgramPath, err := os.Getwd() //获取执行路径
			if err != nil {
				fmt.Printf("%s.ParseConfig:  获取当前执行路径出错。错误信息:%s\n", pack, err)
				return fmt.Errorf("%s.ParseConfig:  获取当前执行路径出错。错误信息:%s\n", pack, err)
			}
			pathArgs := strings.Split(path, "/")
			// 从执行路径开始逐级遍历到插件路径，当目录不存在就直接创建
			path = ProgramPath // 是绝对路径
			for _, DirName := range pathArgs {
				path = path + "/" + DirName
				if ok, _ := Utils.HasDir(path); !ok {
					//路径不存在
					Utils.CreateDir(path) //创建配置文件夹
				}
			}
			// 创建配置文件
			configFile, _ := os.OpenFile(path+"/config.yaml", os.O_CREATE|os.O_RDWR, 0777)
			_, _ = configFile.WriteString(defaultFile) // 将默认的配置写入文件
			_ = configFile.Close()
			// 新建配置文件后仍然因为某种原因出错
			if err = config.ReadInConfig(); err != nil {
				fmt.Printf("%s.ParseConfig:  [%s]配置文件读取出错。错误信息:%s\n", pack, pack, err)
				return fmt.Errorf("%s.ParseConfig:  [%s]配置文件读取出错。错误信息:%s\n", pack, pack, err)
			}
		} else {
			fmt.Printf("%s.ParseConfig:  [%s]配置文件读取出错。错误信息:%s\n", pack, pack, err)
			return fmt.Errorf("%s.ParseConfig:  [%s]配置文件读取出错。错误信息:%s\n", pack, pack, err)
		}
	}
	return nil
}
