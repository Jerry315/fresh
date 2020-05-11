package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
	"test/fresh/models"
)

type UserController struct {
	beego.Controller
}

// 展示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

// 用户注册
func (this *UserController) HandleReg() {
	// 1.获取数据
	name := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	fmt.Println("email: ", email)

	// 2.校验数据
	if name == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整，请重新注册"
		this.TplName = "register.html"
		return
	}

	if pwd != cpwd {
		this.Data["errmsg"] = "两次密码不一致"
		this.TplName = "register.html"
		return
	}

	// 邮箱的正则匹配
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[A-Za-z0-9_-]+(\\.[A-Za-z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不匹配"
		this.TplName = "register.html"
		return
	}
	// 3.处理数据

	o := orm.NewOrm()
	var user models.User
	user.Name = name
	//err := o.Read(&user)
	//if err == nil{
	//
	//}
	user.PassWord = pwd
	user.Email = email
	_, err := o.Insert(&user)
	if err != nil {
		fmt.Printf("插入数据失败：%v\n", err)
		this.Data["errmsg"] = "用户已存在"
		this.TplName = "register.html"
		return
	}

	// email三方秘钥：fjrplzdomgyqbbee
	emailConfig := `{"username":"710899905@qq.com","password":"fjrplzdomgyqbbee","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "710899905@qq.com"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	emailConn.Text = "http://127.0.0.1:8080/active?id=" + strconv.Itoa(user.Id)
	err = emailConn.Send()
	if err != nil {
		fmt.Printf("邮件发送失败：%#v\n", err)
	}
	// 4.返回视图
	this.Ctx.WriteString("注册成功，请去注册邮箱中激活用户")
}

func (this *UserController) ShowActive() {
	uid, _ := this.GetInt("id")
	var user models.User
	user.Id = uid

	o := orm.NewOrm()
	err := o.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "register.html"
		return
	}
	user.Active = true
	o.Update(&user)
	this.Redirect("/user/login", 302)
}

func (this *UserController) ShowLogin() {
	name := this.Ctx.GetCookie("UserName")
	if name != "" {
		nameStr, _ := base64.StdEncoding.DecodeString(name)
		this.Data["UserName"] = string(nameStr)
		this.Data["checked"] = "checked"
		//this.Redirect("/", 302)
		//return
	} else {
		this.Data["UserName"] = ""
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}

func (this *UserController) HandleLogin() {
	// 1.获取数据
	name := this.GetString("username")
	pwd := this.GetString("pwd")
	// 2.校验数据
	if name == "" && pwd == "" {
		this.Data["errmsg"] = "用户名或密码不能为空"
		this.TplName = "login.html"
		return
	}
	// 3.处理数据
	var user models.User
	user.Name = name
	o := orm.NewOrm()
	err := o.Read(&user, "name")
	if err != nil {
		fmt.Printf("用户名不存在: %v\n", err)
		this.Data["errmsg"] = "用户名或密码错"
		this.TplName = "login.html"
		return
	}

	if user.PassWord != pwd {
		fmt.Printf("用户名不存在")
		this.Data["errmsg"] = "用户名或密码错"
		this.TplName = "login.html"
		return
	}
	remember := this.GetString("remember")
	nameStr := base64.StdEncoding.EncodeToString([]byte(name))
	if remember == "on" {
		this.Ctx.SetCookie("UserName", nameStr, 24*3600*30)
	} else {
		this.Ctx.SetCookie("UserName", nameStr, 1)
	}
	//fmt.Println("remember: ",remember)
	this.SetSession("UserName", name)
	// 4.返回视图
	this.Redirect("/", 302)
}

func (this *UserController) HandleLogout() {
	this.DelSession("UserName")
	this.Redirect("/login", 302)
}
