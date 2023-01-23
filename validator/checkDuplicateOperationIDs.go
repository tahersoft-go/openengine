package validator

import "fmt"

func (v *openApiValidator) CheckDuplicateOperationIDs() *openApiValidator {
	allMatches := FindAllByRegex(v.Doc, REGEXP_OPERATION_ID)
	oprationIds := []string{}
	for _, match := range allMatches {
		if len(match) != 2 {

			break
		}
		oprationIds = append(oprationIds, match[1])
	}
	v.OperationIds = oprationIds
	dupValues, isValid := HasSliceDuplicateString(oprationIds)
	v.IsValidOperationIds = isValid
	for _, dupValue := range dupValues {
		BuildError(&v.Errors, fmt.Sprintf("OperationId with value of %s is duplicated", dupValue))
	}
	return v
}
