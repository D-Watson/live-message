package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/D-Watson/live-safety/log"
	"live-message/dbs"
	"live-message/entity"
	"live-message/mq"
	"live-message/utils"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client
	offline  chan *Client
	// Unregister requests from clients.
	unregister chan *Client

	UserHub map[string]*Client

	recieve chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 1024),
		register:   make(chan *Client, 1024),
		offline:    make(chan *Client, 1024),
		unregister: make(chan *Client, 1024),
		UserHub:    make(map[string]*Client, 1024),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.registerClient(ctx, client)

		case client := <-h.unregister:
			delete(h.UserHub, client.userId)

		case c := <-h.offline:
			//todo:用户离线消息存储
			fmt.Println(c)

		case msg := <-h.broadcast:
			deliverMsg(ctx, msg)
		}
	}
}

func (h *Hub) registerClient(ctx context.Context, client *Client) {
	h.UserHub[client.userId] = client
	ip := utils.GetLocalIP()
	err := dbs.SetUserServer(client.userId, ip)
	if err != nil {
		log.Errorf(ctx, "[redis] set ip err=", err)
		return
	}
}

func deliverMsg(ctx context.Context, msg []byte) error {
	data := &entity.Data{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Errorf(ctx, "[json] unmarshal err=", err)
		return err
	}
	user := data.ToUserId
	topic, err := dbs.GetUserServer(user)
	if err != nil {
		log.Errorf(ctx, "[redis] get err=", err)
		return err
	}
	p := mq.NewProducer(topic)
	err = p.Produce(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
