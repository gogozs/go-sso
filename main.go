package main

import (
	"go-weixin/pkg/log"
	"go-weixin/service/api"
)

func main() {
	server := api.StartServer()
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
