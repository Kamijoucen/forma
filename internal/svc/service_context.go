// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"forma/internal/config"
	"forma/internal/ent"
	"forma/internal/ent/migrate"
	"forma/internal/middleware"
	"strings"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	Ent            *ent.Client
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	d := initDB(c)

	// 初始化响应处理器
	initHandler()

	return &ServiceContext{
		Config:         c,
		Ent:            d,
		AuthMiddleware: middleware.NewAuthMiddleware(c).Handle,
	}
}

func initDB(c config.Config) *ent.Client {
	driver := strings.TrimSpace(c.DB.Driver)
	dsn := strings.TrimSpace(c.DB.DSN)
	if driver == "" {
		driver = "postgres"
	}

	if dsn == "" {
		err := fmt.Errorf("db dsn is required")
		logx.Errorf("init db failed: %v", err)
		panic(err)
	}

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		logx.Errorf("init db failed, driver=%s err=%v", driver, err)
		panic(err)
	}
	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true)); err != nil {
		panic(err)
	}

	return client
}
