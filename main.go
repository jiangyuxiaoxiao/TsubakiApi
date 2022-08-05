package main

import (
	"TsubakiApi/Cmd"
	"TsubakiApi/Config"
	"TsubakiApi/Log"
)

func main() {
	Log.InitLogger() //日志初始化
	Config.Parse()   // 配置初始化
	tsubaki.Run()    // 服务器运行
}
