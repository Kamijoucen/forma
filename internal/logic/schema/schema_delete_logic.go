// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/ent"
	"forma/internal/ent/fielddef"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaDeleteLogic {
	return &SchemaDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaDeleteLogic) SchemaDelete(req *types.SchemaDeleteReq) error {
	return util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		// 查找 Schema
		sd, err := tx.SchemaDef.
			Query().
			Where(schemadef.NameEQ(req.Name)).
			Only(l.ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errorx.ErrNotFound
			}
			return err
		}

		// 级联删除关联的 FieldDef
		if _, err := tx.FieldDef.
			Delete().
			Where(fielddef.HasSchemaDefWith(schemadef.IDEQ(sd.ID))).
			Exec(l.ctx); err != nil {
			return err
		}

		// 删除 SchemaDef
		return tx.SchemaDef.DeleteOne(sd).Exec(l.ctx)
	})
}
