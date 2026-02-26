// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"
	"fmt"

	"forma/internal/ent"
	"forma/internal/errorx"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppCreateLogic {
	return &AppCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppCreateLogic) AppCreate(req *types.AppCreateReq) (resp *types.AppCreateResp, err error) {
	a, err := l.svcCtx.Ent.App.Create().
		SetCode(req.Code).
		SetName(req.Name).
		SetDescription(req.Description).
		Save(l.ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errorx.ErrAppAlreadyExists
		}
		return nil, err
	}
	return &types.AppCreateResp{
		Id: fmt.Sprintf("%d", a.ID),
	}, nil
}
