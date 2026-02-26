// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/errorx"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppUpdateLogic {
	return &AppUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppUpdateLogic) AppUpdate(req *types.AppUpdateReq) error {
	a, err := l.svcCtx.Ent.App.Query().
		Where(entApp.CodeEQ(req.Code)).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errorx.ErrAppNotFound
		}
		return err
	}

	update := l.svcCtx.Ent.App.UpdateOne(a)
	if req.Name != "" {
		update.SetName(req.Name)
	}
	if req.Description != "" {
		update.SetDescription(req.Description)
	}
	return update.Exec(l.ctx)
}
