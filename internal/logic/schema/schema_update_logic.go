// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaUpdateLogic {
	return &SchemaUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaUpdateLogic) SchemaUpdate(req *types.SchemaUpdateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
