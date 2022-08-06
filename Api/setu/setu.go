package setu

import "github.com/gin-gonic/gin"

var Setu *gin.RouterGroup

func init() {

}
func Run() {
	Setu.GET("/live", live)
}

func live(context *gin.Context) {

}
