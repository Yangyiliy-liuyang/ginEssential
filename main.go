package main

import (
	"ginEssential/comment"
	"ginEssential/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	_ = comment.GetDB()
	r := gin.Default()
	r = routers.CollectRoute(r)
	err := r.Run()
	if err != nil {
		return
	} // 监听并在 0.0.0.0:8080 上启动服务
}
