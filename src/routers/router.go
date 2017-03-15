package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	/**
         正则路由
	 */
	//*全匹配方式 //匹配 /api/list/file/api.json :splat=file/api.json
	/*beego.Router("/api/list*//*", &controllers.MainController{}, "*:ListFood")
	beego.Router("/api/create*//*", &controllers.MainController{}, "post:CreateFood")
	beego.Router("/api/update*//*", &controllers.MainController{}, "put:UpdateFood")
	beego.Router("/api/delete*//*", &controllers.MainController{}, "delete:DeleteFood")*/


	//APIS
	ns := beego.NewNamespace("/api",
		//过滤器，判断是否登录
		beego.NSCond(func(ctx *context.Context) bool {
			mobile := ctx.GetCookie("mobile")
			beegosessionID := ctx.GetCookie("beegosessionID")
			store := ctx.Input.CruSession

			beego.Info("---------过滤器，判断是否登录", mobile, " Seesion:", store.Get("mobile"), " state", store.SessionID())
			if mobile == "" || store.SessionID() != beegosessionID {
				ctx.Redirect(302, "/tologin")
				return false
			}
			return true
		}),

		beego.NSNamespace("/post",
			//CRUD Create(创建)、Read(读取)、Update(更新)和Delete(删除)
			beego.NSNamespace("/create",
				// /api/post/create/node/
				beego.NSRouter("/*", &controllers.MainController{}, "post:CreateFood"),
			),
			beego.NSNamespace("/read",
				beego.NSRouter("/*", &controllers.MainController{}, "*:ListFood"),
			),
			beego.NSNamespace("/update",
				beego.NSRouter("/*", &controllers.MainController{}, "post:UpdateFood"),
			),
			beego.NSNamespace("/delete",
				beego.NSRouter("/*", &controllers.MainController{}, "post:DeleteFood"),
			)),
	)

	beego.AddNamespace(ns)

	//固定路由
	//路由就是我们最常用的路由方式，一个固定的路由，一个控制器，然后根据用户请求方法不同请求控制器中对应的方法，典型的 RESTful 方式
	beego.Router("/tologin", &controllers.MainController{}, "*:ToLogin")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/register", &controllers.MainController{}, "*:Register")
	/*beego.Router("/postList", &controllers.MainController{}, "*:GetPostList")
	beego.Router("/favortOp", &controllers.MainController{}, "*:FavortOp")
	beego.Router("/deletePost", &controllers.MainController{}, "*:DeletePost")
	beego.Router("/deleteComment", &controllers.MainController{}, "*:DeleteComment")
	beego.Router("/addComment", &controllers.MainController{}, "*:AddComment")
	beego.Router("/addPost", &controllers.MainController{}, "*:AddPost")
	beego.Router("/addUrlPost", &controllers.MainController{}, "*:AddUrlPost")
	beego.Router("/addVideoPost", &controllers.MainController{}, "*:AddVideoPost")*/

	/*//开启session


	Filter()*/
}
func Filter() {

	var FilterUser = func(ctx *context.Context) {
		_, ok := ctx.Input.Session("mobile").(string)
		beego.Info("授权操作:", ok, ctx.Request.RequestURI)
		if !ok && ctx.Request.RequestURI != "/myajax/*" {
			ctx.Redirect(302, "/tologin")
		}
	}

	//授权操作
	beego.InsertFilter("/myajax/*", beego.BeforeRouter, FilterUser)

}