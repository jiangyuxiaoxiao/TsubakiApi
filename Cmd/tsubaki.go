// Package Tsubaki: 主运行文件
package Tsubaki

import (
	"TsubakiApi/Api/setu"
	"TsubakiApi/Api/voice"
	"TsubakiApi/Api/zao"
	"TsubakiApi/Config"
	"TsubakiApi/Log"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Run() {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)
	// gin.SetMode(gin.ReleaseMode) // 设置发布环境
	// 注册枣子路由
	zao.Zao = router.Group("/zao")
	zao.Run()
	// 注册涩图路由
	setu.Setu = router.Group("/setu")
	setu.Run()
	// 注册拟声路由
	voice.Voice = router.Group("/voice")
	voice.Run()
	// 路由启动
	port := Config.Server.Port //端口号
	if port < 0 || port > 65535 {
		Log.Logger.Errorf("Cmd.tsubaki: 端口号不合法。端口号范围0~65535")
		panic("Cmd.tsubaki: 路由格式不合法。")
	}
	portString := ":" + strconv.Itoa(port)
	err := router.Run(portString)
	if err != nil {
		Log.Logger.Errorf("	Cmd.tsubaki: 路由启动失败。")
		panic("	Cmd.tsubaki: 路由启动失败。")
	}

}
