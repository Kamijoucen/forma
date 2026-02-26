// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"
	"fmt"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type EntityCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityCreateLogic {
	return &EntityCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityCreateLogic) EntityCreate(req *types.EntityCreateReq) (resp *types.EntityCreateResp, err error) {
	// 查询 Schema 及其字段定义（按 app 过滤）
	sd, err := l.svcCtx.Ent.SchemaDef.
		Query().
		Where(
			schemadef.NameEQ(req.SchemaName),
			schemadef.HasAppWith(entApp.CodeEQ(req.AppCode)),
		).
		WithFieldDefs().
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorx.NewBizErrorf(errorx.CodeNotFound, "Schema %s 不存在", req.SchemaName)
		}
		return nil, err
	}

	// 校验字段值
	defMap, err := service.ValidateEntityFields(sd.Edges.FieldDefs, req.Fields)
	if err != nil {
		return nil, err
	}

	var recordID int
	if err := util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		// 创建 EntityRecord
		record, err := tx.EntityRecord.Create().
			SetSchemaDef(sd).
			Save(l.ctx)
		if err != nil {
			return err
		}
		recordID = record.ID

		// 批量创建字段值
		creates := lo.Map(req.Fields, func(fv *types.FieldValueInput, _ int) *ent.EntityFieldValueCreate {
			return tx.EntityFieldValue.Create().
				SetValue(fv.Value).
				SetEntityRecord(record).
				SetFieldDef(defMap[fv.Name])
		})
		if _, err := tx.EntityFieldValue.CreateBulk(creates...).Save(l.ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &types.EntityCreateResp{
		Id: fmt.Sprintf("%d", recordID),
	}, nil
}
