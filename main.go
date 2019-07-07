package main

import (
	"go-weixin/pkg/log"
	"go-weixin/service/api"
	"os"
)



func main() {
	server, err := api.StartServer()
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	defer server.Close()
}
