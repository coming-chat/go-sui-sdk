package sui_types

type DynamicFieldType struct {
	DynamicField  *EmptyEnum `json:"DynamicField"`
	DynamicObject *EmptyEnum `json:"DynamicObject"`
}

type DynamicFieldName struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}
