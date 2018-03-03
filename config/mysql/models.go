package mysql

import (
	"log"
	"strconv"
	"time"
	//"github.com/jinzhu/now"
)

type WcAuthorizationList struct {
	RecordId               int    `json:"record_id"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	NickName               string `json:"nick_name"`
}

type ClientList struct {
	ClientId int64  `json:"client_id"`
	OpenId   string `json:"open_id"`
	NickName string `json:"nick_name"`
	AppId    int    `json:"app_id"`
}

type CustomerStatistics struct {
	Id         int64 `json:"id"`
	TaskId     int64 `json:"task_id"`
	AppId      int   `json:"app_id"`
	UserCount  int   `json:"user_count"`
	UserNum    int   `json:"user_num"`
	CreateTime int64 `json:"create_time"`
}

func (CustomerStatistics) TableName() string {
	return "wc_customer_statistics"
}

func (WcAuthorizationList) TableName() string {
	return "wc_authorization_list"
}

//TODO: 根据AppId获取公众号信息
func GetAppInfo(recordId int) (*WcAuthorizationList, bool) {
	appInfo := &WcAuthorizationList{}
	err := db.Select("record_id, authorizer_appid, authorizer_access_token, nick_name").
		Where("record_id = ? and status = 1", recordId).
		First(&appInfo).Error
	if err != nil {
		log.Printf("select getAppInfo err: %v", err)
	}
	if appInfo.RecordId == 0 {
		return appInfo, false
	}
	return appInfo, true
}

//TODO: 获取用户列表
func GetUserList(recordId int) ([]ClientList, bool) {
	var list []ClientList
	//nTime := time.Now().Unix()
	//statTime := nTime - (86400 * 2)
	num := strconv.Itoa(recordId)
	err := db.Table("wc_client" + num).
		Select("client_id, open_id, nick_name, app_id").
			Where("open_id = ?", "ol_EGvw_V3rXYILgc7QEOVVBrxwg").Find(&list).Error
		//Where("update_time between ? and ?", statTime, now.EndOfHour().Unix()).Find(&list).Error
	if err != nil {
		log.Printf("select UserList err : %v", err)
		return nil, false
	}
	return list, true
}
//TODO: 保存任务信息
func SaveTask(AppId int, TaskId int64, count int, num int) (err error) {
	task := &CustomerStatistics{
		AppId:     AppId,
		TaskId:    TaskId,
		UserCount: count,
		UserNum:   num,
		CreateTime: time.Now().Unix(),
	}
	if err := db.Create(&task).Error; err != nil {
		log.Printf("create task record err : %v", err)
		return err
	}
	return nil
}
