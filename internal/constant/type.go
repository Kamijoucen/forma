package constant

const (
	FieldTypeString  = "string"
	FieldTypeNumber  = "number"
	FieldTypeBoolean = "boolean"
	FieldTypeDate    = "date"
	FieldTypeText    = "text"
	FieldTypeEnum    = "enum"
	FieldTypeJSON    = "json"
	FieldTypeArray   = "array"
)

func GetFieldTypes() []string {
	return []string{
		FieldTypeString,
		FieldTypeNumber,
		FieldTypeBoolean,
		FieldTypeDate,
		FieldTypeText,
		FieldTypeEnum,
		FieldTypeJSON,
		FieldTypeArray,
	}
}
