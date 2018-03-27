package main

import (
	_ "chat/routers"

	"github.com/astaxie/beego"
)

const (
	APP_VER = "0.1.1.0227"
)

func main() {
	beego.Info(beego.BConfig.AppName, APP_VER)

	beego.Run()
}
