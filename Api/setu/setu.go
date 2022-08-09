package setu

import (
	"TsubakiApi/Utils"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/big"
)

var Setu *gin.RouterGroup
var FlagLoadConfigError bool = false       //配置文件读取出错标志
var LiveDirError bool = false              //live文件夹读取出错标志
var LiveSetus []string = make([]string, 0) //live文件夹下的所有涩图路径string

func init() {
	// 加载setu插件配置
	err := LoadConfig()
	if err != nil {
		FlagLoadConfigError = true
	}
	// setu/live相关初始化
	ok, err := Utils.HasDir(LiveConfig.AbsPath)
	if !ok || err != nil {
		LiveDirError = true
	} else {
		liveWalkFunc(LiveConfig.AbsPath)
	}

}

func Run() {
	// 配置文件读取出错
	if FlagLoadConfigError {
		return
	}
	// live路由
	if LiveDirError {
		fmt.Printf("setu/live: 加载涩图文件夹出错，清检查配置文件！\n")
	} else {
		Setu.GET("/live", live)
	}

}

func live(context *gin.Context) {
	fileNum := len(LiveSetus)
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(fileNum))) //获取真随机数
	file := LiveSetus[randNum.Int64()]
	//_, name := filepath.Split(file)
	//_ = context.BindJSON(gin.H{"name": name})
	context.File(file)

}

//liveWalkFunc 递归遍历文件夹
func liveWalkFunc(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		// 如果是文件夹
		if file.IsDir() {
			liveWalkFunc(path + "/" + file.Name())
		} else {
			// 如果不是文件夹且是图片
			if Utils.IsImage(file.Name()) {
				// 添加到色批[]string
				LiveSetus = append(LiveSetus, path+"/"+file.Name())
			}
		}
	}
}
