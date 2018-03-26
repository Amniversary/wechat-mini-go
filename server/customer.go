package server

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/Amniversary/wechat-mini-go/controllers"
	"github.com/Amniversary/wechat-mini-go/config/mysql"
	"github.com/Amniversary/wechat-mini-go/config"
)


//TODO: 客服消息
func (s *Server) CustomerMsg(w http.ResponseWriter, r *http.Request) {
	req := &controllers.Customer{}
	rsp := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, rsp)
	}()
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("json decode err: %v", err)
		rsp.Msg = config.ErrMsg
		return
	}
	s.work.Index[req.TaskId+000+req.AppId] = &controllers.Count{
		Success: 0,
		Failed:  0,
		Total:   0,
	}
	auth, ok := mysql.GetAppInfo(req.AppId)
	if !ok {
		log.Printf("CustomerMsg getAppInfo AppId:[%d].", req.AppId)
		rsp.Msg = config.ErrMsg
		return
	}
	list, ok := mysql.GetUserList(auth.RecordId)
	if !ok {
		log.Printf("CustomerMsg getUserList AppId:[%d].", auth.RecordId)
		rsp.Msg = config.ErrMsg
		return
	}
	s.work.Index[req.TaskId+000+req.AppId].Total = len(list)
	//s.work.Task <- list
	go func() {
		for _, v := range list {
			Client := s.work.NewUsers(v, req)
			s.work.Client <- Client
		}
	}()
	rsp.Code = config.RESPONSE_OK
}