// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/ent"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/service"
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
	sd, err := l.svcCtx.Ent.SchemaDef.
		Query().
		Where(schemadef.NameEQ(req.Name)).
		WithFieldDefs().
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorx.ErrNotFound
		}
		return nil, err
	}
	return service.ToSchemaDetailResp(sd), nil
}
