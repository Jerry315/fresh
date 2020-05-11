package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"test/fresh/controllers"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, CheckLogin)
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/active", &controllers.UserController{}, "get:ShowActive")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/user/logout",&controllers.UserController{},"get:HandleLogout")
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
}

var CheckLogin = func(ctx *context.Context) {
	userName := ctx.Input.Session("UserName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
