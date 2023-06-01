package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) { //ctx: context, 上下文，可以理解为一个请求所相关的所有信息
		//1.获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		//2.数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}
		//如何没有传入名称，则随机生成一个10位字符串
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Println(name, telephone, password)
		//3.判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
			return
		}
		//4.创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			password:  password,
		}
		db.Create(&newUser)
		//5.返回结果
		ctx.JSON(200, gin.H{
			"message": "用户注册成功",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
func InitDB() *gorm.DB {
	username := "root"       //账号
	password := "12345"      //密码
	host := "127.0.0.1"      //数据库地址，可以是Ip或者域名
	port := 3306             //数据库端口
	Dbname := "ginessential" //数据库名
	timeout := "10s"         //连接超时，10秒

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 连接成功
	fmt.Println(db)
	db.AutoMigrate(&User{})
	return db
}
func RandomString(n int) string {
	var letters = []byte("assadsasadfsadsaGdsfsFFSWLKDKSKDAXDF")
	result := make([]byte, n)
	//使用时间作为种子值，然后生成不同系列的随机数
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
