// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"
	"strconv"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/ent/entityrecord"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EntityDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityDetailLogic {
	return &EntityDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityDetailLogic) EntityDetail(req *types.EntityDetailReq) (resp *types.EntityDetailResp, err error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, errorx.NewBizError(errorx.CodeInvalidParam, "ID格式不正确")
	}

	record, err := l.svcCtx.Ent.EntityRecord.
		Query().
		Where(
			entityrecord.IDEQ(id),
			entityrecord.HasSchemaDefWith(
				schemadef.NameEQ(req.SchemaName),
				schemadef.HasAppWith(entApp.CodeEQ(req.AppCode)),
			),
		).
		WithFieldValues(func(q *ent.EntityFieldValueQuery) {
			q.WithFieldDef()
		}).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorx.ErrNotFound
		}
		return nil, err
	}

	return service.ToEntityDetailResp(record, req.SchemaName), nil
}
