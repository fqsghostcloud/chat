package controllers

import (
	"container/list"
	"time"

	"chat/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// init chat room
func init() {
	go chatroom()
}

func newEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(userName string, ws *websocket.Conn) {
	comeinChatterCh <- Chatter{Name: userName, Conn: ws}
}

func Leave(user string) {
	exitChatterCh <- user
}

type Chatter struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	comeinChatterCh = make(chan Chatter, 10)
	// Channel for exit users.
	exitChatterCh = make(chan string, 10)
	// Send events here to commonInfoCh them.
	commonInfoCh = make(chan models.Event, 10)

	chatterLists = list.New()
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case chatter := <-comeinChatterCh:
			if !isUserExist(chatterLists, chatter.Name) {
				chatterLists.PushBack(chatter) // Add user to the end of list.
				// Publish a JOIN event.
				commonInfoCh <- newEvent(models.EVENT_JOIN, chatter.Name, "")
				beego.Info("New user:", chatter.Name, ";WebSocket:", chatter.Conn != nil)
			} else {
				beego.Info("Old user:", chatter.Name, ";WebSocket:", chatter.Conn != nil)
			}
		case event := <-commonInfoCh:

			broadcastWebSocket(event)
			// models.AddEvent(event) 从events list 获取消息历史记录

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-exitChatterCh:
			for sub := chatterLists.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Chatter).Name == unsub {
					chatterLists.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Chatter).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					commonInfoCh <- newEvent(models.EVENT_LEAVE, unsub, "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Chatter).Name == user {
			return true
		}
	}
	return false
}
