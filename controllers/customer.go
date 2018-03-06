package controllers

import (
	"net/http"
	"log"
	"encoding/json"
	"sync"
	"context"

	"github.com/Amniversary/wechat-mini-go/config/mysql"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/custom"
	"github.com/Amniversary/wechat-mini-go/config"
)

type Worker struct {
	Client chan *Customer
	Lock   sync.Mutex

	index  map[int64]*Count
	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

type Count struct {
	Success int
	Failed  int
	Total   int
}

type Customer struct {
	AppId   int64  `json:"app_id"`
	KeyWord string `json:"key_word"`
	TaskId  int64  `json:"task_id"`
	MsgData []Msg  `json:"msg_data"`

	OpenId   string `json:"open_id, omitempty"`
	NickName string `json:"nick_name, omitempty"`
	ClientId int64  `json:"client_id, omitempty"`
}

type Msg struct {
	MsgType int              `json:"msg_type"`
	Content string           `json:"content, omitempty"`
	List    []custom.Article `json:"list, omitempty"`
	MediaId string           `json:"media_id, omitempty"`
}

//TODO: 客服消息
func (w *Worker) CustomerMsg(writer http.ResponseWriter, request *http.Request) {
	req := &Customer{}
	rsp := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(writer, http.StatusOK, rsp)
	}()
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		log.Printf("json decode err: %v", err)
		rsp.Msg = config.ErrMsg
		return
	}
	w.index[req.TaskId+000+req.AppId] = &Count{
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
	w.index[req.TaskId+000+req.AppId].Total = len(list)
	go func() {
		for _, v := range list {
			Client := NewUsers(v, req)
			w.Client <- Client
		}
	}()
	rsp.Code = config.RESPONSE_OK
}

//TODO: 发送客服消息
func (w *Worker) SendCustomer(msg *Customer) {
	auth, ok := mysql.GetAppInfo(msg.AppId)
	if !ok {
		return
	}

	var err error
	for _, v := range msg.MsgData {
		var rst interface{}
		switch v.MsgType {
		case 0:
			rst = custom.NewText(msg.OpenId, v.Content, "")
		case 1:
			rst = custom.NewNews(msg.OpenId, v.List, "")
		case 2:
			rst = custom.NewImage(msg.OpenId, v.MediaId, "")
		case 3:
			rst = custom.NewVoice(msg.OpenId, v.MediaId, "")
		}
		err = Send(rst, auth)
	}
	if err != nil {
		w.Lock.Lock()
		w.index[msg.TaskId+000+msg.AppId].Failed ++
		w.Lock.Unlock()
		log.Printf("send customer err: [%v], %v", msg.NickName, err)
	} else {
		w.Lock.Lock()
		w.index[msg.TaskId+000+msg.AppId].Success ++
		w.Lock.Unlock()
	}
	log.Printf("taskId: [%v], appId: [%v], success: [%v], failed: [%v], totel: [%v], userName: [%v], openId: [%v]",
		msg.TaskId,
		msg.AppId,
		w.index[msg.TaskId+000+msg.AppId].Success,
		w.index[msg.TaskId+000+msg.AppId].Failed,
		w.index[msg.TaskId+000+msg.AppId].Total,
		msg.NickName,
		msg.OpenId,
	)
	count := w.index[msg.TaskId+000+msg.AppId].Success + w.index[msg.TaskId+000+msg.AppId].Failed
	if count >= w.index[msg.TaskId+000+msg.AppId].Total {
		err := mysql.SaveTask(msg.AppId, msg.TaskId, w.index[msg.TaskId+000+msg.AppId].Total, w.index[msg.TaskId+000+msg.AppId].Success)
		if err != nil {
			log.Printf("saveTask err: %v", err)
			return
		}
		log.Printf("customer task is over taskId: [%v], appId: [%v], AppName: [%v]", msg.TaskId, msg.AppId, auth.NickName)
	}
	return
	//text := custom.NewText(msg.OpenId,  "测试测试", "")
	//log.Printf("text : %v", text)
}

func NewUsers(v mysql.ClientList, req *Customer) *Customer {
	Client := new(Customer)
	Client.OpenId = v.OpenId
	Client.NickName = v.NickName
	Client.ClientId = v.ClientId
	Client.AppId = v.AppId
	Client.TaskId = req.TaskId
	Client.KeyWord = req.KeyWord
	Client.MsgData = req.MsgData
	return Client
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
