// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EntityUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityUpdateLogic {
	return &EntityUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityUpdateLogic) EntityUpdate(req *types.EntityUpdateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
