// Package Tsubaki: 主运行文件
package Tsubaki

import (
	"TsubakiApi/Api/zao"
	"TsubakiApi/Log"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	zao.Zao = router.Group("/zao") // 注册枣子路由
	zao.Run()
	err := router.Run()
	if err != nil {
		Log.Logger.Errorf("路由启动失败。")
	}
	_ = router.Run()
}
