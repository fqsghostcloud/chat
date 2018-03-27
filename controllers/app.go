package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

// AppController handles the welcome screen that allows user to pick a technology and username.
type AppController struct {
	baseController
}

// Get implemented Get() method for AppController.
func (this *AppController) Get() {
	fmt.Println("hahahaha")
	this.TplName = "welcome.html"
}

// Join method handles POST requests for AppController.
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
	case "longpolling":
		this.Redirect("/lp?uname="+uname, 302)
	case "websocket":
		this.Redirect("/ws?uname="+uname, 302)
	default:
		this.Redirect("/", 302)
	}

	// Usually put return after redirect.
	return
}
