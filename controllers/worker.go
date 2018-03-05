package controllers

import (
	"log"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/Amniversary/wechat-mini-go/config/mysql"
)

const (
	DefaultWorker = 1024
	incompleteURL = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="

	ErrCodeOK                 = 0
	ErrCodeInvalidCredential  = 40001 // access_token 过期错误码
	ErrCodeAccessTokenExpired = 42001 // access_token 过期错误码(maybe!!!)
	ErrCodeUserTimeOut        = 45015 // 用户响应超时
)

type Error struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (w *Worker) initWorker() {
	log.Printf("WeChat init workers.")
	for i := 0; i < DefaultWorker; i++ {
		w.wg.Add(1)
		go w.runWorkers()
	}
}

func (w *Worker) runWorkers() {
	for {
		select {
		case msg := <-w.Client:
			w.SendCustomer(msg)
		case <-w.ctx.Done():
			log.Printf("ctx. Done")
			w.wg.Done()
			return
		}
	}
}

func NewClient() *Worker {
	w := &Worker{
		Client: make(chan *Customer),
		index:  make(map[int64]*Count),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.initWorker()
	return w
}

func Send(msg interface{}, auth *mysql.WcAuthorizationList) (err error) {
	isRefresh := false //TODO: 重试发送
RETRY:
	var result Error
	Url := incompleteURL + auth.AuthorizerAccessToken //TODO: 请求客服消息 url
	client := &http.Client{}
	buf := bytes.NewBuffer([]byte{}) //TODO: 使用json.Marshal 会使特殊字符unicode
	Json := json.NewEncoder(buf)
	Json.SetEscapeHTML(false)
	Json.Encode(msg)
	if err != nil {
		log.Printf("json encode err : %v", err)
		return
	}
	req, err := http.NewRequest("POST", Url, buf)
	if err != nil {
		log.Printf("http new request err : %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http do request err : %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code err: %s", resp.StatusCode)
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil realAll err: %v", err)
		return
	}
	if err := json.Unmarshal(rspBody, &result); err != nil { //TODO: json.decode
		log.Printf("json unmarshal err : %v", err)
		return err
	}
	switch result.ErrCode {
	case ErrCodeOK:
		return
	case ErrCodeInvalidCredential, ErrCodeAccessTokenExpired: //TODO: token过期处理(异常)
		if !isRefresh {
			isRefresh = true
			Auth, ok := mysql.GetAppInfo(auth.RecordId)
			if !ok {
				return
			}
			auth.AuthorizerAccessToken = Auth.AuthorizerAccessToken
			goto RETRY
		}
		return fmt.Errorf("ErrCode: [%v], ErrMsg: [%v]", result.ErrCode, result.ErrMsg)
	case ErrCodeUserTimeOut: //TODO:	用户交互超时
		return fmt.Errorf("ErrCodeUserTimeOut ErrCode : [%v], ErrMsg: [%v]", result.ErrCode, result.ErrMsg)
	default:
		return fmt.Errorf("ErrCode: [%v], ErrMsg: [%v]", result.ErrCode, result.ErrMsg)
	}

	return
}
