package consts

import (
	"time"
)

//ws的常量

const (
	WriteTimeout   = 1 * time.Second
	MAXMessageSize = 1024
	ReadTimeout    = 1 * time.Second
	PongPeriod     = (ReadTimeout * 9) / 10
)
