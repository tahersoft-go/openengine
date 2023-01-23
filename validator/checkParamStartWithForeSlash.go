package validator

import "fmt"

func (v *openApiValidator) CheckParamStartWithForeSlash() *openApiValidator {
	for path := range v.YamlDoc.Paths {
		if path[0] != '/' {
			BuildError(&v.Errors, fmt.Sprintf("Path %s does not start with /", path))
		}
	}
	return v
}
