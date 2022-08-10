package voice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

var Voice *gin.RouterGroup
var FlagLoadConfigError bool = false //配置文件读取出错标志
var FlagYuzuStartError bool = false  //yuzu初始化启动标志

var YuzuIn io.WriteCloser
var YuzuOut io.ReadCloser
var WriteBuff []byte = make([]byte, 4096)
var lock sync.Mutex
var yuzuCmd *exec.Cmd

func init() {
	// 加载voice插件配置
	err := LoadConfig()
	if err != nil {
		FlagLoadConfigError = true
	}
	// atri相关初始化
	// yuzu相关初始化 进行交互式控制台处理
	yuzuCmd = exec.Command("cmd", "/C", YuzuConfig.Vits)
	YuzuIn, _ = yuzuCmd.StdinPipe()
	yuzuCmd.Stdout = os.Stdout
	err = yuzuCmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	_, _ = io.WriteString(YuzuIn, YuzuConfig.ModulePath+"\n")
	_, _ = io.WriteString(YuzuIn, YuzuConfig.Config+"\n")
}

func Run() {
	// 配置文件读取出错
	if FlagLoadConfigError {
		return
	}
	// atri路由
	Voice.Handle("GET", "/atri", atri)
	// yuzu路由
	if !FlagYuzuStartError {
		Voice.Handle("GET", "/yuzu", yuzu)
	}
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

func yuzu(context *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	_, _ = io.WriteString(YuzuIn, "t\n")
	// 获取文本
	text, _ := context.GetQuery("text")
	text = text + "\n"
	_, _ = io.WriteString(YuzuIn, text)
	// 获取选择的人物
	id, _ := context.GetQuery("id")
	_, _ = io.WriteString(YuzuIn, id+"\n")
	// 获取存放路径
	path := YuzuConfig.Output
	_, _ = io.WriteString(YuzuIn, path+"/output.wav\n")
	// 再次循环
	_, _ = io.WriteString(YuzuIn, "y\n")
	// 发送请求
	time.Sleep(2000 * time.Millisecond)
	context.File(YuzuConfig.Output + "/output.wav")
}
