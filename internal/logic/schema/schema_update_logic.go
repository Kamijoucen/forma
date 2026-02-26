// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/constant"
	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/ent/fielddef"
	"forma/internal/ent/schemadef"
	"forma/internal/errorx"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaUpdateLogic {
	return &SchemaUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaUpdateLogic) SchemaUpdate(req *types.SchemaUpdateReq) error {
	if err := service.ValidateSchemaFields(req.Fields); err != nil {
		return err
	}

	return util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		// 查找 Schema（按 app 过滤）
		sd, err := tx.SchemaDef.
			Query().
			Where(
				schemadef.NameEQ(req.Name),
				schemadef.HasAppWith(entApp.CodeEQ(req.AppCode)),
			).
			Only(l.ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errorx.ErrNotFound
			}
			return err
		}

		// 更新 SchemaDef 可变字段
		if err := tx.SchemaDef.UpdateOne(sd).
			SetDescription(req.Description).
			Exec(l.ctx); err != nil {
			return err
		}

		// 查出关联的所有 FieldDef，按 name 建索引
		existingFields, err := tx.FieldDef.
			Query().
			Where(fielddef.HasSchemaDefWith(schemadef.IDEQ(sd.ID))).
			All(l.ctx)
		if err != nil {
			return err
		}
		fieldMap := lo.SliceToMap(existingFields, func(fd *ent.FieldDef) (string, *ent.FieldDef) {
			return fd.Name, fd
		})

		// 遍历请求字段，按 name 匹配更新可变属性
		for _, reqField := range req.Fields {
			existing, ok := fieldMap[reqField.Name]
			if !ok {
				return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 不存在，不支持新增字段", reqField.Name)
			}
			// type 不可变，校验一致性
			if string(existing.Type) != reqField.Type {
				return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的类型不可修改", reqField.Name)
			}

			// 枚举去重
			enumValues := reqField.EnumValues
			if reqField.Type == constant.FieldTypeEnum {
				enumValues = lo.Uniq(enumValues)
			}

			if err := tx.FieldDef.UpdateOne(existing).
				SetRequired(reqField.Required).
				SetMaxLength(reqField.MaxLength).
				SetMinLength(reqField.MinLength).
				SetEnumValues(enumValues).
				SetDescription(reqField.Description).
				Exec(l.ctx); err != nil {
				return err
			}
		}
		return nil
	})
}
