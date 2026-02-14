// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"forma/internal/config"
	"forma/internal/ent"
)

type ServiceContext struct {
	Config config.Config
	Ent    *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	d := initDB(c)

	// 初始化响应处理器
	initHandler()

	return &ServiceContext{
		Config: c,
		Ent:    d,
	}
}

func initDB(c config.Config) *ent.Client {
	// TODO
	return nil
}
