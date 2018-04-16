package main

import (
	"github.com/Amniversary/wechat-mini-go/server"
	"github.com/Amniversary/wechat-mini-go/config"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server.NewServer(config.NewConfig()).Run()
}
