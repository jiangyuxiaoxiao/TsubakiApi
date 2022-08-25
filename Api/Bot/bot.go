package Bot

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"sync"
)

var Bot *gin.RouterGroup

type Messages struct {
	MessageLock sync.Mutex
	Messages    []Message
}

var GlobalMessages Messages

type Message struct {
	QQ      string //QQ号
	Text    string //接收者
	IsGroup bool   //是群聊
}

// 初始化函数
func init() {
	GlobalMessages.Messages = make([]Message, 0)
}

func Run() {
	// 接收Bot信息路由路由
	Bot.Handle("GET", "/receive", receive)
	// 发送Bot信息路由
	Bot.Handle("GET", "/send", send)
}

func receive(context *gin.Context) {
	text, _ := context.GetQuery("text")
	qq, _ := context.GetQuery("qq")
	isGroup, _ := context.GetQuery("isGroup")
	var ok bool
	if isGroup == "true" {
		ok = true
	} else {
		ok = false
	}
	GlobalMessages.MessageLock.Lock()
	defer GlobalMessages.MessageLock.Unlock()
	msg := Message{
		QQ:      qq,
		Text:    text,
		IsGroup: ok,
	}
	// 将接收到的消息添加到队列中
	GlobalMessages.Messages = append(GlobalMessages.Messages, msg)
	context.JSON(200, "OK")
}

func send(context *gin.Context) {
	GlobalMessages.MessageLock.Lock()
	defer GlobalMessages.MessageLock.Unlock()
	jsonByte, _ := json.Marshal(GlobalMessages.Messages)
	jsonString := string(jsonByte)
	GlobalMessages.Messages = GlobalMessages.Messages[0:0] //清空切片
	context.JSON(200, jsonString)
}
