// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/constant"
	"forma/internal/ent"
	"forma/internal/ent/fielddef"
	"forma/internal/errorx"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaCreateLogic {
	return &SchemaCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaCreateLogic) SchemaCreate(req *types.SchemaCreateReq) error {

	orm := l.svcCtx.Ent
	ctx := l.ctx

	if err := validateSchema(req); err != nil {
		return err
	}

	if err := util.WithTx(ctx, orm, func(tx *ent.Tx) error {
		return l.Do(ctx, tx, req)
	}); err != nil {
		return err
	}
	return nil
}

func (l *SchemaCreateLogic) Do(ctx context.Context, tx *ent.Tx, req *types.SchemaCreateReq) error {

	sd, err := tx.SchemaDef.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		Save(ctx)
	if err != nil {
		return err
	}

	convert := func(r *types.FieldDef) (*ent.FieldDefCreate, error) {

		// 枚举去重
		enumValues := r.EnumValues
		if r.Type == constant.FieldTypeEnum {
			enumValues = lo.Uniq(enumValues)
		}

		create := tx.FieldDef.Create().
			SetName(r.Name).
			SetType(fielddef.Type(r.Type)).
			SetRequired(r.Required).
			SetMaxLength(r.MaxLength).
			SetMinLength(r.MinLength).
			SetEnumValues(enumValues).
			SetDescription(r.Description).
			SetSchemaDef(sd)
		return create, nil
	}

	deflist := lo.Map(req.Fields, func(r types.FieldDef, _ int) *ent.FieldDefCreate {
		if fd, err := convert(&r); err != nil {
			return nil
		} else {
			return fd
		}
	})

	if _, err = tx.FieldDef.CreateBulk(deflist...).Save(ctx); err != nil {
		return err
	}
	return nil
}

func validateSchema(req *types.SchemaCreateReq) error {

	if len(req.Fields) == 0 {
		return errorx.ErrInvalidParam
	}
	for _, f := range req.Fields {
		if err := fielddef.TypeValidator(fielddef.Type(f.Type)); err != nil {
			return err
		}

		if f.Type == constant.FieldTypeEnum && len(f.EnumValues) == 0 {
			return errorx.NewBizError(errorx.CodeInvalidParam, "枚举类型必须提供枚举值")
		}
	}
	return nil
}
