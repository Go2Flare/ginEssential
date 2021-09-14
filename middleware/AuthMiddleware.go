package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oceanlearn.learn/ginessential/common"
	"oceanlearn.learn/ginessential/model"
	"strings"
	"log"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {

		//获取authorization header
		tokenString:=ctx.GetHeader("Authorization")

		//	validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer "){
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, Claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid{
			log.Printf("%v, token = %v", err, token.Valid)
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}

		//验证通过获取claim中的userID
		userId := Claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户不存在
		if userId == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}

		//用户信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}