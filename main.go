package main

import (
	"github.com/gin-gonic/gin"
)

// ginHtml函数
func ginHtml(c *gin.Context) {
	type UserInfo struct {
		UserName string
		Age      int
		PassWord string
	}
	user := UserInfo{"lihan", 32, "123456"}
	c.HTML(200, "index.html", user)

	//c.HTML(200, "index.html", gin.H{"user_name": "lihan", "age": 32, "status": http.StatusOK, "data": gin.H{"id": 1, "name": "lihan"}})
} //gin.H()可以向html传参

func main() {
	router := gin.Default()
	//加载html模板目录下所有模板文件
	//templates目录要与main.go所在目录同级，而不是在main.go所在目录
	router.LoadHTMLGlob("templates/*")

	//golang中，没有相对文件的路径，只有相对项目的路径

	router.Static("/static", "static") //静态文件路径

	router.GET("/html", ginHtml)

	router.Run(":8080")
}
