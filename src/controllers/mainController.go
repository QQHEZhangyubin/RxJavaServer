package controllers

import (
	"github.com/astaxie/beego"
	"./been"
	"../models"
	"time"
)

type MainController struct {
	beego.Controller
}

const (
	IamgePath = "static/upload/images/"
	VideoPath = "static/upload/videos/"
	VideoImgPath = "static/upload/videos/img/"
)

/*
RESTful Controller 路由

在介绍这三种 beego 的路由实现之前先介绍 RESTful
，我们知道 RESTful 是一种目前 API 开发中广泛采用的形式，
beego 默认就是支持这样的请求方法，也就是用户 Get 请求就执行 Get 方法，
Post 请求就执行 Post 方法。因此默认的路由是这样 RESTful 的请求方式。
 */
func (this *MainController)ListFood() {
	splat := this.Ctx.Input.Param(":splat")
	switch splat {
	case "postlist":
		this.GetPostList();
		break
	case "one":

		break


	}

}
func (this *MainController)CreateFood() {
	splat := this.Ctx.Input.Param(":splat")
	switch splat {
	case "post":
		optype := this.Input().Get("type")
		switch optype {
		case "1":
			this.AddPost()
			break
		case "2":
			this.AddUrlPost()
			break
		case "3":
			this.AddVideoPost()
			break
		}
		break
	case "comment":
		this.AddComment()
		break
	case "favort":
		this.AddFavort()
		break
	}
}
func (this *MainController)UpdateFood() {

}
func (this *MainController)DeleteFood() {
	splat := this.Ctx.Input.Param(":splat")
	switch splat {
	case "post":
		this.DeletePost()
		break
	case "comment":
		this.DeleteComment()
		break
	case "favort":
		this.DeleteFavort()
		break
	}
}
/**
 登录
 http://localhost:8080/login?mobile=18820792655&pwd=123456
 */
func (this *MainController)Login() {
	user := &been.UserFrom{}
	err := this.ParseForm(user)
	if err == nil {
		login, exist := models.Login(user)
		if exist {
			beego.Info("登录成功", exist)
			this.SaveSession(login)
			ReturnSuccess(&(this.Controller), login)

		} else {
			beego.Info("帐号或密码错误", exist)
			ReturnError(&(this.Controller), "帐号或密码错误", FAIL)
		}
	} else {
		ReturnError(&(this.Controller), "", ERROR)
	}
}
/**
  注册
  http://localhost:8080/register?mobile=18820792655&pwd=123456
 */
func (this *MainController)Register() {
	user := &been.UserFrom{}
	err := this.ParseForm(user)
	if err == nil {
		u, err := models.Register(user)
		if err == nil {
			ReturnSuccess(&(this.Controller), u)
		} else {
			ReturnError(&(this.Controller), err.Error(), FAIL)
		}
	} else {
		ReturnError(&(this.Controller), "", ERROR)
	}
}

func (this *MainController)ToLogin() {
	ReturnError(&(this.Controller), "请登录", TOLOGIN)
}

func (this *MainController)Logout() {
	this.Ctx.SetCookie("mobile", "", 0, "/")
	this.Ctx.SetCookie("pwd", "", 0, "/")
	this.DestroySession()
	ReturnSuccess(&(this.Controller), "")

}


/**
保存Session
 */
func (this *MainController)SaveSession(u *models.User) {
	//30分钟有效
	maxTime := 30 * time.Minute
	this.Ctx.SetCookie("mobile", u.Mobile, maxTime, "/")
	this.Ctx.SetCookie("pwd", u.Pwd, maxTime, "/")
	this.SetSession("mobile", u.Mobile)

	beego.Info("-----------SaveSession  ----- ", u.Mobile)

}


