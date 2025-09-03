package main

import (
	"context"

	cf "github.com/D-Watson/live-safety/conf"
	"live-message/controller"
	"live-message/dbs"
)

func main() {
	ctx := context.Background()
	err := cf.ParseConfig(ctx)
	if err != nil {
		return
	}
	dbs.InitMysql()
	dbs.InitRedis(ctx)
	controller.InitWsServer(ctx)
}
