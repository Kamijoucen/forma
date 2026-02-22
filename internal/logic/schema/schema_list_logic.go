// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/service"
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
	list, err := l.svcCtx.Ent.SchemaDef.
		Query().
		WithFieldDefs().
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*types.SchemaDetailResp, 0, len(list))
	for _, sd := range list {
		items = append(items, service.ToSchemaDetailResp(sd))
	}
	return &types.SchemaListResp{
		Total: int64(len(list)),
		List:  items,
	}, nil
}
