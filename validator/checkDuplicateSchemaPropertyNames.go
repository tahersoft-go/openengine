package validator

import "fmt"

func (v *openApiValidator) CheckDuplicateSchemaPropertyNames() *openApiValidator {
	for schemaName, schema := range v.YamlDoc.Components.Schemas {
		mapStringInterface, err := ToMapStringInterface(schema.Properties)
		if err != nil {
			BuildError(&v.Errors, err.Error())
			continue
		}
		keys := GetMapKeys(mapStringInterface)
		dupValues, hasDuplicateValue := HasSliceDuplicateString(keys)
		if !hasDuplicateValue {
			continue
		}
		for _, dupValue := range dupValues {
			BuildError(&v.Errors, fmt.Sprintf("Schema %s has duplicated property name %s", schemaName, dupValue))
		}
	}
	return v
}
