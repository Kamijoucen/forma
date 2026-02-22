// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"

	"forma/internal/ent"
	"forma/internal/ent/entityrecord"
	"forma/internal/ent/schemadef"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/samber/lo"
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
	query := l.svcCtx.Ent.EntityRecord.
		Query().
		Where(entityrecord.HasSchemaDefWith(schemadef.NameEQ(req.SchemaName)))

	// 查询总数
	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	// 分页查询
	records, err := query.
		WithFieldValues().
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order(entityrecord.ByID()).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := lo.Map(records, func(r *ent.EntityRecord, _ int) *types.EntityDetailResp {
		return service.ToEntityDetailResp(r, req.SchemaName)
	})

	return &types.EntityListResp{
		Total: int64(total),
		List:  list,
	}, nil
}
