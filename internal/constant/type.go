package constant

const (
	FieldTypeString  = "string"
	FieldTypeNumber  = "number"
	FieldTypeBoolean = "boolean"
	FieldTypeDate    = "date"
)

func GetFieldTypes() []string {
	return []string{
		FieldTypeString,
		FieldTypeNumber,
		FieldTypeBoolean,
		FieldTypeDate,
	}
}
