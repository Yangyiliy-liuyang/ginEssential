package main

import (
	"ginEssential/comment"
	"ginEssential/routers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()
	_ = comment.GetDB()
	r := gin.Default()
	r = routers.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	//err := r.Run()
	//if err != nil {
	//	return
	//} // 监听并在 0.0.0.0:8080 上启动服务
}
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {

	}
}
