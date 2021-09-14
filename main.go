package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"oceanlearn.learn/ginessential/common"
	"os"
)


func main() {
	//导入配置文件
	InitConfig()
	//初始化数据库
	common.InitDB()
	//初始化框架
	r := gin.Default()
	//路由连接
	r = CollectRoute(r)
	//viper读取配置
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}

	panic(r.Run())
	//r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil{
		panic(err)
	}
}

