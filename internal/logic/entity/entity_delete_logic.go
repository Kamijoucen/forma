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
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type EntityDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityDeleteLogic {
	return &EntityDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityDeleteLogic) EntityDelete(req *types.EntityDeleteReq) error {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return errorx.NewBizError(errorx.CodeInvalidParam, "ID格式不正确")
	}

	return util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		// 验证 EntityRecord 存在且属于该 Schema
		record, err := tx.EntityRecord.
			Query().
			Where(
				entityrecord.IDEQ(id),
				entityrecord.HasSchemaDefWith(schemadef.NameEQ(req.SchemaName)),
			).
			Only(l.ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errorx.ErrNotFound
			}
			return err
		}

		// 级联删除关联的字段值
		if _, err := tx.EntityFieldValue.
			Delete().
			Where(entityfieldvalue.HasEntityRecordWith(entityrecord.IDEQ(record.ID))).
			Exec(l.ctx); err != nil {
			return err
		}

		// 删除 EntityRecord
		return tx.EntityRecord.DeleteOne(record).Exec(l.ctx)
	})
}
