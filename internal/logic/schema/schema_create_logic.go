// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"context"

	"forma/internal/constant"
	"forma/internal/ent"
	"forma/internal/ent/fielddef"
	"forma/internal/service"
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

	if err := service.ValidateSchemaFields(req.Fields); err != nil {
		return err
	}

	// 查找 App
	app, err := service.QueryAppByCode(l.ctx, l.svcCtx.Ent, req.AppCode)
	if err != nil {
		return err
	}

	if err := util.WithTx(l.ctx, l.svcCtx.Ent, func(tx *ent.Tx) error {
		return l.Do(tx, req, app)
	}); err != nil {
		return err
	}
	return nil
}

func (l *SchemaCreateLogic) Do(tx *ent.Tx, req *types.SchemaCreateReq, app *ent.App) error {

	sd, err := tx.SchemaDef.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetApp(app).
		Save(l.ctx)
	if err != nil {
		return err
	}

	convertFn := func(r *types.FieldDef) *ent.FieldDefCreate {
		// 枚举去重
		enumValues := r.EnumValues
		if r.Type == constant.FieldTypeEnum {
			enumValues = lo.Uniq(enumValues)
		}
		return tx.FieldDef.Create().
			SetName(r.Name).
			SetType(fielddef.Type(r.Type)).
			SetRequired(r.Required).
			SetMaxLength(r.MaxLength).
			SetMinLength(r.MinLength).
			SetEnumValues(enumValues).
			SetDescription(r.Description).
			SetSchemaDef(sd)
	}

	defList := lo.Map(req.Fields, func(fd *types.FieldDef, _ int) *ent.FieldDefCreate {
		return convertFn(fd)
	})

	if _, err := tx.FieldDef.CreateBulk(defList...).Save(l.ctx); err != nil {
		return err
	}
	return nil
}
