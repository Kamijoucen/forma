// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaCreateLogic {
	return &SchemaCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaCreateLogic) SchemaCreate(req *types.SchemaCreateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
