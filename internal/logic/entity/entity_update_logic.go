// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"
	"strconv"

	"forma/internal/ent"
	"forma/internal/ent/entityfieldvalue"
	"forma/internal/ent/entityrecord"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/samber/lo"
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
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return errorx.NewBizError(errorx.CodeInvalidParam, "ID格式不正确")
	}

	// 查询 Schema 及其字段定义
	sd, err := l.svcCtx.Ent.SchemaDef.
		Query().
		Where(schemadef.NameEQ(req.SchemaName)).
		WithFieldDefs().
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errorx.NewBizErrorf(errorx.CodeNotFound, "Schema %s 不存在", req.SchemaName)
		}
		return err
	}

	// 校验字段值
	defMap, err := service.ValidateEntityFields(sd.Edges.FieldDefs, req.Fields)
	if err != nil {
		return err
	}

	return util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		// 验证 EntityRecord 存在且属于该 Schema
		record, err := tx.EntityRecord.
			Query().
			Where(
				entityrecord.IDEQ(id),
				entityrecord.HasSchemaDefWith(schemadef.IDEQ(sd.ID)),
			).
			Only(l.ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errorx.ErrNotFound
			}
			return err
		}

		// 删除旧的字段值
		if _, err := tx.EntityFieldValue.
			Delete().
			Where(entityfieldvalue.HasEntityRecordWith(entityrecord.IDEQ(record.ID))).
			Exec(l.ctx); err != nil {
			return err
		}

		// 批量创建新的字段值
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
	})
}
