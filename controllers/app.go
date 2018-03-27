package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

type AppController struct {
	baseController
}

func (this *AppController) Get() {
	fmt.Println("hahahaha")
	this.TplName = "welcome.html"
}

func (this *AppController) Join() {
	// Get form value.
	uname := this.GetString("uname")
	tech := this.GetString("tech")

	// Check valid.
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	switch tech {
	case "websocket":
		this.Redirect("/ws?uname="+uname, 302)
	default:
		this.Redirect("/", 302)
	}

	// Usually put return after redirect.
	return
}
