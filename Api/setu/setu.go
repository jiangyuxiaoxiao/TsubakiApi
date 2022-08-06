package setu

import (
	"github.com/gin-gonic/gin"
)

var Setu *gin.RouterGroup
var FlagLoadConfigError bool //配置文件读取出错标志

func init() {
	FlagLoadConfigError = false
	err := LoadConfig() //加载setu插件配置
	if err != nil {
		FlagLoadConfigError = true
	}
}

func Run() {
	// 配置文件读取出错
	if FlagLoadConfigError {
		return
	}
	Setu.GET("/live", live)
}

func live(context *gin.Context) {

}
