package controllers

import (
	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

type AppController struct {
	baseController
}

func (this *AppController) Get() {
	this.TplName = "welcome.html"
}

func (this *AppController) Join() {
	// Get form value.
	uname := this.GetString("uname")

	// Check valid.
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.Redirect("/ws?uname="+uname, 302)

	// Usually put return after redirect.
	return
}
