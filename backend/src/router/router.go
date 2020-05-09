package router

import (
	"controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine{
	//初始化路由引擎对象
	r:=gin.Default()
	// CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.Default())

	//定义路由，以及对应动作处理函数
	r.GET("/ping", controller.Ping)
	r.GET("/category-tree",controller.CategoryTree)
	r.POST("/category",controller.CategoryAdd)

	//返回路由引擎对象
	return r
}
