package main

import "go-weixin/service/api"

func main() {
	//server := api.StartServer()
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Error(err)
	//}
	router := api.InitRouter()
	router.Run("localhost:8003")
}
