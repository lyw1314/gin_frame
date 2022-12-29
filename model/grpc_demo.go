package model

import (
	"context"
	pb "gin_frame/grpc_pb/api_data"
)

type Auth struct {
	AppId        string `json:"app_id"`
	AppSecret    string `json:"app_secret"`
	DeveloperId  string `json:"developer_id"`
	CreatTime    int64  `json:"creat_time"`
	ModifiedTime int64  `json:"modified_time"`
}

func (a *Auth) GetList(condition map[string]interface{}) (list []Auth, err error) {
	//err = DbHandle["master"].Where(condition).Find(&list).Error
	return
}

func (a *Auth) GetOne(appId string) (err error) {
	//a.AppId = "test123"
	//a.AppSecret = "test456"
	//return
	//
	//err = DbHandle["master"].Where("app_id = ?", appId).Last(&a).Error
	//return
	// gRPC调用
	var one *pb.AuthGetOneResponse
	one, err = ApiDataClient.AuthGetOne(context.Background(), &pb.AuthGetOneRequest{AppId: appId})
	if err != nil {
		return
	}
	a.AppId = one.AppId
	a.AppSecret = one.AppSecret
	a.CreatTime = one.CreatTime
	a.DeveloperId = one.DeveloperId
	return
}
