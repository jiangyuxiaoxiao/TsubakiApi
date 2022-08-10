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
	yuzuCmd = exec.Command(YuzuConfig.Vits)
	YuzuIn, _ = yuzuCmd.StdinPipe()
	yuzuCmd.Stdout = os.Stdout
	err = yuzuCmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	WriteBuff = []byte(YuzuConfig.ModulePath + "\n")
	_, err = YuzuIn.Write(WriteBuff)
	if err != nil {
		FlagYuzuStartError = true
	}
	if FlagYuzuStartError {
		fmt.Printf("voice/yuzu 控制台启动失败。错误信息%s\n", err)
	}
	WriteBuff = []byte(YuzuConfig.Config + "\n")
	_, err = YuzuIn.Write(WriteBuff)
	if err != nil {
		FlagYuzuStartError = true
	}
	if FlagYuzuStartError {
		fmt.Printf("voice/yuzu 控制台启动失败。错误信息%s\n", err)
	}
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
	typeString := "t\n"
	WriteBuff = []byte(typeString)
	_, _ = YuzuIn.Write(WriteBuff)
	fmt.Println(WriteBuff)
	// 获取文本
	text, _ := context.GetQuery("text")
	text = text
	WriteBuff = []byte(text + "\n")
	_, _ = YuzuIn.Write(WriteBuff)
	fmt.Println(WriteBuff)
	// 获取选择的人物
	id, _ := context.GetQuery("id")
	WriteBuff = []byte(id + "\n")
	_, _ = YuzuIn.Write(WriteBuff)
	fmt.Println(WriteBuff)
	// 获取存放路径
	path := YuzuConfig.Output + "/output.wav\n"
	WriteBuff = []byte(path)
	_, _ = YuzuIn.Write(WriteBuff)
	fmt.Println(WriteBuff)
	// 再次循环
	WriteBuff = []byte("y\n")
	_, _ = YuzuIn.Write(WriteBuff)
	fmt.Println(WriteBuff)

	/*
		text, _ := context.GetQuery("text")
		text = text + "\n"
		yuzuCmd.Stdin.Read([]byte(text))
		id, _ := context.GetQuery("id")
		id = id + "\n"
		yuzuCmd.Stdin.Read([]byte(id))
		path := YuzuConfig.Output + "/output.wav\n"
		path = path + "\n"
		myrune, _ := utf8.DecodeRuneInString(path)
		utf8.EncodeRune(WriteBuff, myrune)
		yuzuCmd.Stdin.Read(WriteBuff)
		content := "y\n"
		yuzuCmd.Stdin.Read([]byte(content))
	*/
	// 发送请求
	time.Sleep(2000 * time.Millisecond)
	context.File(YuzuConfig.Output + "/output.wav")
}
