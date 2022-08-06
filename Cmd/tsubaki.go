// Package Tsubaki: 主运行文件
package Tsubaki

import (
	"TsubakiApi/Api/setu"
	"TsubakiApi/Api/zao"
	"TsubakiApi/Log"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	// 注册枣子路由
	zao.Zao = router.Group("/zao")
	zao.Run()
	// 注册涩图路由
	setu.Setu = router.Group("/setu")
	err := router.Run(":1210")
	//
	if err != nil {
		Log.Logger.Errorf("路由启动失败。")
	}

}
