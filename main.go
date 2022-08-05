package main

import (
	"TsubakiApi/Cmd"
	"TsubakiApi/Config"
	"TsubakiApi/Log"
)

func main() {
	Config.Parse()   // 配置初始化
	Log.InitLogger() //日志初始化
	Tsubaki.Run()    // 服务器运行
}
