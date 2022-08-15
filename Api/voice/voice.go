package voice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var Voice *gin.RouterGroup
var FlagLoadConfigError bool = false //配置文件读取出错标志
var FlagYuzuStartError bool = false  //yuzu初始化启动标志

var YuzuIn io.WriteCloser
var YuzuOut io.ReadCloser
var WriteBuff []byte = make([]byte, 4096)
var lock sync.Mutex
var YuzuCmd *exec.Cmd
var YuzuLock []sync.Mutex

func init() {
	// 加载voice插件配置
	err := LoadConfig()
	if err != nil {
		FlagLoadConfigError = true
	}
	// atri相关初始化
	// yuzu相关初始化 进行交互式控制台处理
	YuzuLock = make([]sync.Mutex, YuzuConfig.MaxConcurrent) //柚子锁
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
	var lockNumber int // 锁编号
	for i, _ := range YuzuLock {
		ok := YuzuLock[i].TryLock()
		if ok {
			lockNumber = i
			break
		}
	}
	defer YuzuLock[lockNumber].Unlock()
	YuzuCmd = exec.Command("cmd", "/C", "python MoeGoe.py")
	YuzuCmd.Dir = YuzuConfig.GoeMoePythonPath
	YuzuIn, _ = YuzuCmd.StdinPipe()
	err := YuzuCmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	_, _ = io.WriteString(YuzuIn, YuzuConfig.ModulePath+"\n")
	_, _ = io.WriteString(YuzuIn, YuzuConfig.Config+"\n")
	_, _ = io.WriteString(YuzuIn, "t\n")
	// 获取文本
	text, _ := context.GetQuery("text")
	text = text + "\n"
	fileName := YuzuConfig.StringFile + "/" + strconv.Itoa(lockNumber) + ".txt" //文件名 与锁对应
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	file.WriteString(text)
	file.Close()
	_, _ = io.WriteString(YuzuIn, fileName+"\n")
	// 获取选择的人物
	id, _ := context.GetQuery("id")
	_, _ = io.WriteString(YuzuIn, id+"\n")
	// 获取存放路径
	path := YuzuConfig.Output
	path = path + "/" + strconv.Itoa(lockNumber) + ".wav" //文件名 与锁对应
	_, _ = io.WriteString(YuzuIn, path+"\n")
	// 再次循环
	_, _ = io.WriteString(YuzuIn, "n\n")
	// 发送请求
	YuzuCmd.Wait()
	context.File(path)
}
