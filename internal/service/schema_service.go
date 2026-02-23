package service

import (
	"time"

	"forma/internal/constant"
	"forma/internal/ent"
	"forma/internal/ent/fielddef"
	"forma/internal/errorx"
	"forma/internal/types"
)

// ValidateSchemaFields 校验字段定义列表：非空、名称唯一、类型合法、枚举类型必须有值
func ValidateSchemaFields(fields []*types.FieldDef) error {
	if len(fields) == 0 {
		return errorx.ErrInvalidParam
	}
	nameSet := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if _, exists := nameSet[f.Name]; exists {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段名 %s 重复", f.Name)
		}
		nameSet[f.Name] = struct{}{}
		if err := fielddef.TypeValidator(fielddef.Type(f.Type)); err != nil {
			return errorx.NewBizError(errorx.CodeInvalidParam, err.Error())
		}
		if f.Type == constant.FieldTypeEnum && len(f.EnumValues) == 0 {
			return errorx.NewBizError(errorx.CodeInvalidParam, "枚举类型必须提供枚举值")
		}
	}
	return nil
}

// ToSchemaDetailResp 将 Ent SchemaDef（需 eager load FieldDefs）转为 API 响应
func ToSchemaDetailResp(sd *ent.SchemaDef) *types.SchemaDetailResp {
	fields := make([]*types.FieldDef, 0, len(sd.Edges.FieldDefs))
	for _, fd := range sd.Edges.FieldDefs {
		fields = append(fields, ToFieldDef(fd))
	}
	return &types.SchemaDetailResp{
		Name:        sd.Name,
		Description: sd.Description,
		Fields:      fields,
		CreatedAt:   sd.CreateTime.Format(time.DateTime),
		UpdatedAt:   sd.UpdateTime.Format(time.DateTime),
	}
}

// ToFieldDef 将 Ent FieldDef 实体转为 API 类型
func ToFieldDef(fd *ent.FieldDef) *types.FieldDef {
	return &types.FieldDef{
		Name:        fd.Name,
		Type:        string(fd.Type),
		Required:    fd.Required,
		MaxLength:   fd.MaxLength,
		MinLength:   fd.MinLength,
		EnumValues:  fd.EnumValues,
		Description: fd.Description,
	}
}
