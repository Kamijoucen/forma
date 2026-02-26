// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"forma/internal/ent"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type AppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppListLogic {
	return &AppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppListLogic) AppList() (resp *types.AppListResp, err error) {
	list, err := l.svcCtx.Ent.App.Query().All(l.ctx)
	if err != nil {
		return nil, err
	}

	items := lo.Map(list, func(a *ent.App, _ int) *types.AppDetailResp {
		return service.ToAppDetailResp(a)
	})
	return &types.AppListResp{
		Total: int64(len(list)),
		List:  items,
	}, nil
}
