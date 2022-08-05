package zao

import "github.com/gin-gonic/gin"

var Zao *gin.RouterGroup

func Run() {
	Zao.GET("/hello", hello)
}

func hello(context *gin.Context) {
	context.JSON(200, gin.H{"message": "Hi my name is Zao！"})
}
