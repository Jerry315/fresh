package controllers

import (
	"github.com/astaxie/beego"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndex() {
	user := this.GetSession("UserName")
	if user == nil {
		this.Data["user"] = ""
	} else {
		this.Data["user"] = user
	}

	this.TplName = "index.html"
}
