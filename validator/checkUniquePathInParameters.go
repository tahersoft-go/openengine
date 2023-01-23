package validator

import "fmt"

func (v *openApiValidator) CheckUniquePathInParameters() *openApiValidator {
	mapPathParameters := map[string][]string{}
	for path, operations := range v.YamlDoc.Paths {
		if operations.Get != nil {
			mapPathParameters["get"] = append(mapPathParameters["get"], path)
		}
		if operations.Post != nil {
			mapPathParameters["post"] = append(mapPathParameters["post"], path)
		}
		if operations.Put != nil {
			mapPathParameters["put"] = append(mapPathParameters["put"], path)
		}
		if operations.Patch != nil {
			mapPathParameters["patch"] = append(mapPathParameters["patch"], path)
		}
		if operations.Delete != nil {
			mapPathParameters["delete"] = append(mapPathParameters["delete"], path)
		}
	}
	for method, paths := range mapPathParameters {
		dupValues, hasDuplicateValue := HasSliceDuplicateString(paths)
		if !hasDuplicateValue {
			continue
		}
		for _, dupValue := range dupValues {
			BuildError(&v.Errors, fmt.Sprintf("Method %s has duplicated path %s", method, dupValue))
		}
	}
	return v
}
