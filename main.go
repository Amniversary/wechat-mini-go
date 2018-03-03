package main

import (
	"github.com/Amniversary/wechat-mini-go/server"
	"github.com/Amniversary/wechat-mini-go/config"
)

func main() {
	server.NewServer(config.NewConfig()).Run()
}
