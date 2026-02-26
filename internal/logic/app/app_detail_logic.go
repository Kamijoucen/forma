// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppDetailLogic {
	return &AppDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppDetailLogic) AppDetail(req *types.AppDetailReq) (resp *types.AppDetailResp, err error) {
	a, err := service.QueryAppByCode(l.ctx, l.svcCtx.Ent, req.Code)
	if err != nil {
		return nil, err
	}
	return service.ToAppDetailResp(a), nil
}
