// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaDetailLogic {
	return &SchemaDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaDetailLogic) SchemaDetail(req *types.SchemaDetailReq) (resp *types.SchemaDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
