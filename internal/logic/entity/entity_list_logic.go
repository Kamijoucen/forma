// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EntityListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityListLogic {
	return &EntityListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityListLogic) EntityList(req *types.EntityListReq) (resp *types.EntityListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
