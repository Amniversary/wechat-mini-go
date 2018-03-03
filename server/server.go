package server

import (
	"net/http"
	"log"

	"github.com/Amniversary/wechat-mini-go/config"
	"github.com/Amniversary/wechat-mini-go/config/mysql"
	"github.com/Amniversary/wechat-mini-go/controllers"
)

type ServerBase interface {
	Run()
}

type Server struct {
	cfg  *config.Config
	work *controllers.Worker
}

func NewServer(cfg *config.Config) ServerBase {
	return &Server{cfg: cfg}
}

func (s *Server) init() {
	mysql.NewMysql(s.cfg)
	s.work = controllers.NewClient()
}

func (s *Server) runServer() {
	http.HandleFunc("/rpc", s.rpc)
	log.Printf("ListenServer Port: [%s]", s.cfg.Port)
	http.ListenAndServe(s.cfg.Port, nil)
}

func (s *Server) Run() {
	s.init()
	s.runServer()
}
