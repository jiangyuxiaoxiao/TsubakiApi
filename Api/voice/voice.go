package voice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var Voice *gin.RouterGroup
var FlagLoadConfigError bool = false //配置文件读取出错标志
var FlagYuzuStartError bool = false  //yuzu初始化启动标志

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
	// moegoe进程初始化
	for i := 0; i < YuzuConfig.MaxConcurrent; i++ {
		arg := "python MoeGoe.py" + " " + strconv.Itoa(12100+i)
		YuzuCmd = exec.Command("cmd", "/C", arg)
		YuzuCmd.Dir = YuzuConfig.GoeMoePythonPath
		// YuzuIn, _ = YuzuCmd.StdinPipe()
		err = YuzuCmd.Start()
		if err != nil {
			fmt.Println(err)
		}
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
	// 目前包含柚子社十一女主模型推理
	var ok bool
	var lockNumber int // 锁编号
	for i, _ := range YuzuLock {
		ok = YuzuLock[i].TryLock()
		if ok {
			lockNumber = i
			break
		}
	}
	if !ok {
		context.JSON(404, "")
	}
	defer YuzuLock[lockNumber].Unlock()
	//解析分流 id=0-6 yuzu id=7-11 星巴克 id=？404
	// 获取选择的人物
	id, _ := context.GetQuery("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(404, "")
		return
	}
	switch {
	case 0 <= idNum && idNum <= 6: // 千恋万花
		Path := YuzuConfig.ModulePath
		Config := YuzuConfig.Config
		MoeGoeHandle(context, lockNumber, 0, idNum, Path, Config)
	case 7 <= idNum && idNum <= 11: // 星光咖啡馆
		idNum = idNum - 7
		Path := YuzuConfig.StellaPath
		Config := YuzuConfig.StellaConfig
		MoeGoeHandle(context, lockNumber, 1, idNum, Path, Config)
	case idNum == 12: // 亚托莉
		idNum = 0
		Path := YuzuConfig.AtriPath
		Config := YuzuConfig.AtriConfig
		MoeGoeHandle(context, lockNumber, 2, idNum, Path, Config)
	case 13 <= idNum && idNum <= 20: // 魔女的夜宴等
		idNum = idNum - 13
		Path := YuzuConfig.SabbatPath
		Config := YuzuConfig.SabbatConfig
		MoeGoeHandle(context, lockNumber, 3, idNum, Path, Config)
	case 21 <= idNum && idNum <= 24: // 缘之空
		idNum = idNum - 21
		Path := YuzuConfig.SoraPath
		Config := YuzuConfig.SoraConfig
		MoeGoeHandle(context, lockNumber, 4, idNum, Path, Config)
	case 25 <= idNum && idNum <= 32: // 灵感满溢的甜蜜创想
		idNum = idNum - 25
		Path := YuzuConfig.HamiPath
		Config := YuzuConfig.HamiConfig
		MoeGoeHandle(context, lockNumber, 5, idNum, Path, Config)
	case 33 <= idNum && idNum <= 38: //星空列车与白的旅行
		idNum = idNum - 33
		Path := YuzuConfig.HoshishiroPath
		Config := YuzuConfig.HoshishiroConfig
		MoeGoeHandle(context, lockNumber, 6, idNum, Path, Config)
	case 39 <= idNum && idNum <= 51: //落忆13人模型
		idNum = idNum - 39
		Path := YuzuConfig.Luoyi13Path
		Config := YuzuConfig.Luoyi13Config
		MoeGoeHandle(context, lockNumber, 7, idNum, Path, Config)
	default:
		context.JSON(404, "")
		return
	}

}
func MoeGoeHandle(context *gin.Context, lockNumber int, modelNum int, idNum int, Path string, Config string) {
	// 获取文本
	text, _ := context.GetQuery("text")
	path := YuzuConfig.Output
	path = path + "/" + strconv.Itoa(lockNumber) + ".wav" //文件名 与锁对应
	myurl := fmt.Sprintf("http://127.0.0.1:%d/voice?model_id=%d&chr_id=%d&typ=t&txt=%s&advance=False",
		lockNumber+12100, modelNum, idNum, url.QueryEscape(text))
	// 打开文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		context.JSON(404, "")
		return
	}

	resp, err := http.Get(myurl)
	if err != nil {
		context.JSON(404, "")
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		context.JSON(404, "")
		return
	}
	_ = file.Close()
	_ = resp.Body.Close()
	context.File(path)
}
