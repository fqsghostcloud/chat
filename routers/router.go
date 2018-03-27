package routers

import (
	"chat/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.AppController{})

	beego.Router("/join", &controllers.AppController{}, "post:Join")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

}
