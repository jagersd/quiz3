package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const(
    writeWait = 10 * time.Second
    pongWait = 60* time.Second
    pingPeriod = (pongWait * 9) / 10
    maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

type connection struct{
    ws *websocket.Conn
    send chan []byte
    host bool
}

func (s subscription) readPump() {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
        player,msgType,msgString := ParseMessage(msg)
        if player == quizStates[s.room].Host {
            c.host = true
        }
        playermsg, hostmsg := createResponse(player, msgType, msgString, quizStates[s.room])
		m := message{playermsg, s.room}
		h.sendToPlayers <- m

        x := message{hostmsg, s.room}
        h.sendToHost <- x
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomId := vars["quizSlug"]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump()
}

//returning player -- messagetype -- actual message
func ParseMessage(received []byte) (string,string,string){
    msgString := string(received)
    msgArray := strings.Split(msgString,"|")

    if len(msgArray) != 3{
        return "","",""
    }
    
    return msgArray[0], msgArray[1], msgArray[2]
}
