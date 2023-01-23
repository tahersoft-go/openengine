package validator

import "fmt"

func (v *openApiValidator) CheckUniqueGetParameters() *openApiValidator {
	for path, operations := range v.YamlDoc.Paths {
		if operations.Get == nil {
			continue
		}
		if operations.Get.Parameters == nil {
			continue
		}
		var parameterNames []string
		for _, parameter := range operations.Get.Parameters {
			parameterNames = append(parameterNames, parameter.Name)
		}
		dupValues, hasDuplicateValue := HasSliceDuplicateString(parameterNames)
		if !hasDuplicateValue {
			continue
		}
		for _, dupValue := range dupValues {
			BuildError(&v.Errors, fmt.Sprintf("Path %s has duplicated parameter %s", path, dupValue))
		}
	}
	return v
}
