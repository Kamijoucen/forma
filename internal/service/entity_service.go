package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"forma/internal/ent"
	"forma/internal/ent/fielddef"
	"forma/internal/errorx"
	"forma/internal/types"

	"github.com/samber/lo"
)

// ValidateEntityFields 根据 FieldDef 定义校验 FieldValueInput 列表：必填、长度、枚举值等
// 校验通过后返回 name → *ent.FieldDef 映射，供调用方关联 FieldDef
func ValidateEntityFields(fieldDefs []*ent.FieldDef, fieldValues []*types.FieldValueInput) (map[string]*ent.FieldDef, error) {
	// 构建 FieldDef 索引 name → *ent.FieldDef
	defMap := lo.SliceToMap(fieldDefs, func(fd *ent.FieldDef) (string, *ent.FieldDef) {
		return fd.Name, fd
	})

	// 记录已提供的字段名，用于后续必填检查
	provided := make(map[string]bool, len(fieldValues))

	for _, fv := range fieldValues {
		def, ok := defMap[fv.Name]
		if !ok {
			return nil, errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 未在Schema中定义", fv.Name)
		}

		// 值校验
		if err := validateFieldValue(def, fv); err != nil {
			return nil, err
		}

		provided[fv.Name] = true
	}

	// 必填字段检查
	for _, def := range fieldDefs {
		if def.Required && !provided[def.Name] {
			return nil, errorx.NewBizErrorf(errorx.CodeInvalidParam, "必填字段 %s 未提供", def.Name)
		}
	}

	return defMap, nil
}

// validateFieldValue 根据字段类型校验单个 FieldValueInput
func validateFieldValue(def *ent.FieldDef, fv *types.FieldValueInput) error {
	value := fv.Value

	switch def.Type {
	case fielddef.TypeString, fielddef.TypeText:
		length := utf8.RuneCountInString(value)
		if def.MinLength > 0 && length < def.MinLength {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 长度不能小于 %d", fv.Name, def.MinLength)
		}
		if def.MaxLength > 0 && length > def.MaxLength {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 长度不能大于 %d", fv.Name, def.MaxLength)
		}

	case fielddef.TypeNumber:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值不是合法的数字", fv.Name)
		}

	case fielddef.TypeBoolean:
		if value != "true" && value != "false" {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值必须为 true 或 false", fv.Name)
		}

	case fielddef.TypeDate:
		if _, err := time.Parse(time.DateTime, value); err != nil {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值不是合法的日期格式（%s）", fv.Name, time.DateTime)
		}

	case fielddef.TypeEnum:
		if !lo.Contains(def.EnumValues, value) {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值 %s 不在枚举范围内 %v", fv.Name, value, def.EnumValues)
		}

	case fielddef.TypeJSON:
		if !json.Valid([]byte(value)) {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值不是合法的JSON", fv.Name)
		}

	case fielddef.TypeArray:
		if !json.Valid([]byte(value)) {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值不是合法的JSON", fv.Name)
		}
		var arr []any
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return errorx.NewBizErrorf(errorx.CodeInvalidParam, "字段 %s 的值必须是JSON数组", fv.Name)
		}
	}

	return nil
}

// ToEntityDetailResp 将 Ent EntityRecord（需 eager load FieldValues 及其 FieldDef 边）转为 API 响应
func ToEntityDetailResp(record *ent.EntityRecord, schemaName string) *types.EntityDetailResp {
	fields := lo.Map(record.Edges.FieldValues, func(fv *ent.EntityFieldValue, _ int) *types.FieldValue {
		return &types.FieldValue{
			Name:  fv.Edges.FieldDef.Name,
			Type:  string(fv.Edges.FieldDef.Type),
			Value: fv.Value,
		}
	})
	return &types.EntityDetailResp{
		Id:         fmt.Sprintf("%d", record.ID),
		SchemaName: schemaName,
		Fields:     fields,
		CreatedAt:  record.CreateTime.Format(time.DateTime),
		UpdatedAt:  record.UpdateTime.Format(time.DateTime),
	}
}
