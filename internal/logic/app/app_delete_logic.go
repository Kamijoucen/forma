// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppDeleteLogic {
	return &AppDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppDeleteLogic) AppDelete(req *types.AppDeleteReq) error {
	// 检查 App 下是否有 Schema
	exist, err := l.svcCtx.Ent.SchemaDef.Query().
		Where(schemadef.HasAppWith(entApp.CodeEQ(req.Code))).
		Exist(l.ctx)
	if err != nil {
		return err
	}
	if exist {
		return errorx.ErrAppInUse
	}

	// 按 code 查找并删除
	a, err := l.svcCtx.Ent.App.Query().
		Where(entApp.CodeEQ(req.Code)).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errorx.ErrAppNotFound
		}
		return err
	}
	return l.svcCtx.Ent.App.DeleteOne(a).Exec(l.ctx)
}
