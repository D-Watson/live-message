package controller

import (
	"context"
	"net/http"

	cf "github.com/D-Watson/live-safety/conf"
	"github.com/D-Watson/live-safety/log"
	"github.com/gorilla/mux"
	"live-message/service"
)

func InitWsServer(ctx context.Context) {
	router := mux.NewRouter()
	h := service.NewHub()
	go h.Run(ctx)
	router.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		service.ServeWs(h, writer, request)
	})
	if err := http.ListenAndServe(cf.GlobalConfig.Server.Http.Host, router); err != nil {
		log.Errorf(ctx, "[WS] run server error, err=", err)
		return
	}
}
