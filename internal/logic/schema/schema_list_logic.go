// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaListLogic {
	return &SchemaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaListLogic) SchemaList() (resp *types.SchemaListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
