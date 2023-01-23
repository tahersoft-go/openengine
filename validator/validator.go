package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Errors []error

func (errs Errors) Error() string {
	content := []string{}
	for _, err := range errs {
		content = append(content, err.Error())
	}
	errorContent, err := json.Marshal(content)
	if err != nil {
		return err.Error()
	}
	return string(errorContent)
}

type openApiValidator struct {
	YamlDoc      *apiGen
	Doc          string
	OperationIds []string

	IsValidOperationIds bool

	Errors Errors
}

func (v *openApiValidator) GetErrors() *Errors {
	return &v.Errors
}

func ValidateRaw(rawYamlDoc string) *Errors {
	var yamlDoc apiGen
	if err := yaml.Unmarshal([]byte(rawYamlDoc), &yamlDoc); err != nil {
		if err != nil {
			return &Errors{err}
		}
		return nil
	}

	validator := &openApiValidator{
		YamlDoc: &yamlDoc,
	}

	errors := validator.
		CheckAllRefsExistsInSchema().
		CheckDuplicateOperationIDs().
		CheckDuplicateSchemaPropertyNames().
		CheckParamStartWithForeSlash().
		CheckUniqueGetParameters().
		CheckUniquePathInParameters().
		GetErrors()

	if len(*errors) > 0 {
		return errors
	}
	return nil
}

func ValidateFile(filePath string) *Errors {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return &Errors{err}
	}
	return ValidateRaw(string(bytes))
}

func PrintErrors(errors Errors) {
	fmt.Printf("\n\nResult of OpenApi Validation: [%d error(s) found]\n-----------------\n", len(errors))
	for _, err := range errors {
		log.Println(err.Error(), "\n-----------------")
	}
}
