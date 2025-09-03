package entity

import (
	"time"
)

const (
	MsgTypeText  = "text"
	MsgTypeImage = "image"
	MsgTypeVideo = "video"
	MsgTypeFile  = "file"
	MsgTypeVoice = "voice"
	MsgHeartBeat = "heart"
)

type Data struct {
	FromUserId string `json:"fromUserId"` //发送方userId
	ToUserId   string `json:"toUserId"`   //接收方
	SessionId  string `json:"sessionId"`  //唯一索引，会话id
	Type       string `json:"type"`       //消息类型
	Content    string `json:"content"`    //消息实体
	InfoJson   string `json:"infoJson"`   //可扩展字段
	// 时序与控制
	Timestamp  int64 `json:"timestamp" bson:"timestamp"`   // 客户端时间戳
	ServerTime int64 `json:"serverTime" bson:"serverTime"` // 服务器时间戳
	SeqId      int64 `json:"seqId" bson:"seqId"`           // 会话内顺序ID
	Status     int   `json:"status" bson:"status"`         // 消息状态
	// 数据库相关
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"` // 数据库相关
}

type GroupData struct {
}
