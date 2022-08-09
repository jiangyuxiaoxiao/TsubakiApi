package voice

import (
	"github.com/gin-gonic/gin"
	"os/exec"
)

var Voice *gin.RouterGroup
var FlagLoadConfigError bool = false //配置文件读取出错标志

func init() {
	// 加载voice插件配置
	err := LoadConfig()
	if err != nil {
		FlagLoadConfigError = true
	}
	// atri相关初始化

}

func Run() {
	// 配置文件读取出错
	if FlagLoadConfigError {
		return
	}
	// atri路由
	Voice.Handle("GET", "/atri", atri)
}

func atri(context *gin.Context) {
	text, _ := context.GetQuery("text")
	if text == "" {
		context.JSON(200, gin.H{"status": 400, "error code": "Invalid Params"})
		return
	}
	modulePath := AtriConfig.ModulePath
	outputPath := AtriConfig.OutPut
	exePath := AtriConfig.Tacotron
	noisePath := AtriConfig.NoiseFile
	// 罗马音转音频
	Rome := text
	cmd := exec.Command("./inference.exe", "-m", modulePath, "-t", Rome+".", "-o", outputPath, "-od", noisePath)
	cmd.Dir = exePath
	_ = cmd.Run()
	// 发送音频文件
	context.File(outputPath + ".wav")
}
