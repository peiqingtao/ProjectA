package main

import (
	"config"
	"controller"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"router"
)

func main(){
	//初始化配置
	config.InitConfig()

	switch config.App["APP_MODE"] {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}

	//初始化路由引擎对象
	r:=router.RouterInit()

	//初始化工作
	db, err := controller.InitDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	r.Run(config.App["SERVER_ADDR"]) // listen and serve on 0.0.0.0:8088
}
