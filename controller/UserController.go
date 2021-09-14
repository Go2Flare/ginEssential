package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"oceanlearn.learn/ginessential/common"
	"oceanlearn.learn/ginessential/dto"
	"oceanlearn.learn/ginessential/model"
	"oceanlearn.learn/ginessential/response"
	"oceanlearn.learn/ginessential/utils"
)

func Register(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	fmt.Println(name,telephone,password)
	if len(telephone) != 11 {
		//可以传入42, 也可以传入http的状态常量
		//map[string]interface{}{"code":422, "msg":"手机号必须为11位"}这样写也可以, gin.H是内置的
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"密码不能少于6位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	//如果没有写昵称, 给一个10位的随机字符串
	if len(name) == 0 {
		name = utils.RandomString(10)
	}
	log.Println(name, telephone, password)

	//判断手机号是否存在(查数据库)
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"用户已存在")
		return
	}

	//创建用户
	hassesPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hassesPassword),
	}
	DB.Create(&newUser)

	//返回结果
	ctx.JSON(200, gin.H{
		"code":200,
		"msg": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	DB:=common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		//可以传入42, 也可以传入http的状态常量
		//map[string]interface{}{"code":422, "msg":"手机号必须为11位"}这样写也可以, gin.H是内置的

		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"密码不能少于6位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0{
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil,"用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		response.Fail(ctx,"密码错误", gin.H{"code":400,"msg":"密码错误"})
		ctx.JSON(http.StatusBadRequest, gin.H{"code":400,"msg":"密码错误"})
		//log.Fatal(err)
		return
	}
	//发放token
	token,err := common.ReleaseToken(user)
	if err != nil{

		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500,"msg":"系统异常"})
		log.Printf("token generate error:%v", err)
		return
	}


	//返回结果
	response.Success(ctx, gin.H{"token":token},"登录成功")
}

func Info(ctx *gin.Context) {
	//这里返回的是interface，需要断言
	user, _ := ctx.Get("user")
	//只需要用户名，密码字段

	ctx.JSON(http.StatusOK, gin.H{"code":200,"data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
}

//是否数据库存在手机号
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}