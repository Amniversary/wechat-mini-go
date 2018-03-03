package server

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
)

const (
	AccessName = "WeChatMini"
)

func (s *Server) rpc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("Method not be Post Request [%s]\n", r.Method)
		return
	}
	serverName := r.Header.Get("ServerName")
	if serverName != AccessName {
		log.Printf("rpc service name[%s], not be self [%s] .", serverName, AccessName)
		return
	}
	methodName := r.Header.Get("MethodName")
	start := time.Now()
	defer func() {
		log.Printf("Request MethodName: [%s], Rtime[%v]\n", methodName, time.Now().Sub(start))
	}()

	switch methodName {
	case "customerMsg":
		s.work.CustomerMsg(w, r)
	}
}

// TODO @ 输出Json数据
func EchoJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type,servername,methodname,userid,msgid")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}