package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/D-Watson/live-safety/log"
	"github.com/gorilla/websocket"
	"live-message/consts"
	"live-message/entity"
	kf "live-message/mq"
	"live-message/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  consts.MAXMessageSize,
	WriteBufferSize: consts.MAXMessageSize,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	conn   *websocket.Conn
	userId string
}

// readPump pumps msgs from the ws connection to the hub.
func (c *Client) ReadMessage(ctx context.Context, h *Hub) {
	defer func() {
		// client登出
		h.unregister <- c
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	c.conn.SetReadLimit(consts.MAXMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(consts.ReadTimeout))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorf(ctx, "[WS] read msg error,err=", err)
			break
		}
		h.broadcast <- message
	}
}

func (c *Client) SendMessageByUser(ctx context.Context, h *Hub, msg []byte) error {
	defer func() {
		// client登出
		h.unregister <- c
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	c.conn.SetWriteDeadline(time.Now().Add(consts.WriteTimeout))
	if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Errorf(ctx, "[ws] write msg err=", err)
		return err
	}
	return nil
}

// kafka消费者启动

func OnMessage(ctx context.Context, h *Hub) {
	ip := utils.GetLocalIP()
	consumer := kf.NewConsumer(ip, "")
	err := consumer.ConsumeWithHandler(ctx, func(msg []byte) error {
		data := &entity.Data{}
		err := json.Unmarshal(msg, &data)
		if err != nil {
			log.Errorf(ctx, "[json] unmarshal err=", err)
			return err
		}
		userId := data.ToUserId
		if userId == "" {
			return errors.New("userId is nil")
		}
		c := h.UserHub[userId]
		err = c.SendMessageByUser(ctx, h, msg)
		if err != nil {
			log.Errorf(ctx, "[msg] send err=", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf(ctx, "[onmessage] err=", err)
		return
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf(ctx, "", err)
		return
	}
	//从header获取userId
	userId := r.Header.Get("X-User-Id")
	client := &Client{conn: conn, userId: userId}
	hub.register <- client
	go client.ReadMessage(ctx, hub)
	go OnMessage(ctx, hub)
}
