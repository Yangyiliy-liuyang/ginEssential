package controller

import (
	"ginEssential/comment"
	"ginEssential/dto"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) { //ctx: context, 上下文，可以理解为一个请求所相关的所有信息
	DB := comment.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	//1.获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	//2.数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//如何没有传入名称，则随机生成一个10位字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	//3.判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}
	//4.创建用户
	//加密用户密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "加密密码错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	//发放token
	token, err := comment.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("Token generate error ：%v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
func Login(c *gin.Context) {
	DB := comment.GetDB()
	//获取参数
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号已存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := comment.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("Token generate error ：%v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "用户登录成功")
}
func Info(context *gin.Context) {
	user, _ := context.Get("user")
	context.JSON(http.StatusOK, gin.H{
		"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}
