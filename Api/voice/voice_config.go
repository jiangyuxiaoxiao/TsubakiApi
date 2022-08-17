package voice

import (
	"TsubakiApi/Config"
	_ "embed"
	"fmt"
	"github.com/spf13/viper"
)

//go:embed default_config.yaml
var defaultConfig string

//Atri Atri相关语音配置
type Atri struct {
	ModulePath string `yaml:"ModulePath"` //预训练模型路径
	OutPut     string `yaml:"outPut"`     //输出语音路径
	Tacotron   string `yaml:"Tacotron"`   //Tacotron2路径
	NoiseFile  string `yaml:"NoiseFile"`  //降噪音频路径
}

//Yuzu Yuzu相关语音配置
type Yuzu struct {
	ModulePath       string `yaml:"ModulePath"`       //柚子社模型的绝对路径
	Config           string `yaml:"Config"`           //柚子社配置文件相关路径
	Output           string `yaml:"Output"`           //输出路径
	GoeMoePythonPath string `yaml:"GoeMoePythonPath"` //VGoeMoePython文件路径
	StringFile       string `yaml:"StringFile"`       //缓存日文设置路径
	MaxConcurrent    int    `yaml:"MaxConcurrent"`    //最大并发数
	StellaPath       string `yaml:"StellaPath"`       //星光咖啡馆模型绝对路径
	StellaConfig     string `yaml:"StellaConfig"`     //星光咖啡馆配置文件路径
}

var AtriConfig Atri
var YuzuConfig Yuzu

// LoadConfig 配置加载
func LoadConfig() error {
	config := viper.New()
	config.AddConfigPath("./Api/voice")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	//加载配置文件
	err := Config.PhraseYamlConfig("voice", "/Api/voice", defaultConfig, config)
	if err != nil {
		return err
	} // 直接返回错误信息到上一级
	//加载atri相关配置
	err = config.UnmarshalKey("atri", &AtriConfig)
	if err != nil {
		fmt.Printf("voice/atri Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
		return fmt.Errorf("voice/atri Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
	}
	//加载yuzu相关配置
	err = config.UnmarshalKey("yuzu", &YuzuConfig)
	if err != nil {
		fmt.Printf("voice/yuzu Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
		return fmt.Errorf("voice/yuzu Config 配置加载出错，一般为配置文件格式出错导致。错误信息: %s\n", err)
	}
	//加载配置完成
	return nil
}
