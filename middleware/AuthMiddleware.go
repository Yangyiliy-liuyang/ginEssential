package middleware

import (
	"ginEssential/comment"
	"ginEssential/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware 验证中间键
/**
使用gin的中间件实现用户认证的示例。保护路由和用户验证
主要流程如下:
	1. 获取请求头中的Authorization,得到token字符串
	2. 验证token的格式,必须以Bearer开头
	3. 如果格式错误,返回401 Unauthorized
	4. 调用comment.ParseToken解析JWT token
	5. 如果解析失败或JWT不合法,返回401 Unauthorized
	6. 解析成功后,从JWT的claims中获取用户ID
	7. 根据用户ID从数据库查询用户信息
	8. 如果用户不存在,返回401 Unauthorized
	9. 如果查询到用户,将用户信息写入请求上下文context
	10. 调用next()放行请求
所以,这个中间件实现了:
	1. 从请求中获取JWT token
	2. 解析和验证JWT token
	3. 根据JWT的claims获取用户ID
	4. 根据用户ID查询用户信息
	5. 如果以上任何一步失败,返回401 Unauthorized
	6. 如果全部通过,将用户信息放入上下文,传递给下一个handler
这是一个标准的JWT认证中间件实现,通过解析JWT获取用户信息,并根据查询结果控制请求是否放行。
一些改进的地方:
1. 可以使用更加安全的加密算法替代HS256,如RS256
2. 应有限制解析JWT的时间,避免重放攻击
3. 可以在登录接口中返回Refresh JWT,实现 jwt 刷新
4. 返回的错误信息不应该这么详细,以免暴露系统信息
*/
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取authorization header
		tokenString := context.GetHeader("Authorization")
		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := comment.ParseToken(tokenString)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}
		//验证通过后获取claim中的userid
		userId := claims.UserID
		DB := comment.GetDB()
		var user model.User
		DB.First(&user, userId)
		//用户
		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}
		//用户存在 将user的信息写入上下文
		context.Set("user", user)
		context.Next()
	}
}
