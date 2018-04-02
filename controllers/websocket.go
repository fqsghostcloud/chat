package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"chat/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

// Get method handles GET requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Safe check.
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["UserName"] = uname
}

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Join() {
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
			return
		} else {
			beego.Error("Cannot setup WebSocket connection:", err)
			return
		}
	}

	// Join chat room.
	Join(uname, ws)
	defer Leave(uname)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("\nwebsocket message receive loop erro[%s]\n", err.Error())
			return
		}

		commonInfoCh <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for chatterItem := chatterLists.Front(); chatterItem != nil; chatterItem = chatterItem.Next() {
		// Immediately send event to WebSocket users.
		ws := chatterItem.Value.(Chatter).Conn //断言
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				exitChatterCh <- chatterItem.Value.(Chatter).Name
			}
		}
	}
}
