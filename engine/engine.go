package engine

import (
	"os"
	"regexp"
	"strings"
)

func ExportAPIDocsYaml(dest, content string) error {
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func IsIgnoredFile(filePath string) bool {
	// check in not valid slice
	for _, v := range IGNORED_FILES_TO_PARS {
		if strings.Contains(filePath, v) {
			return true
		}
	}
	return false
}

func ParseStructTagValues(tag string) OpenApiFieldTagValues {
	values := OpenApiFieldTagValues{}
	if tag != "" {
		for _, item := range strings.Split(tag, ";") {
			splitted := strings.Split(item, ":")
			if len(splitted) == 1 {
				values.Required = TerIf(splitted[0] == "required", true, false)
				values.Nullable = TerIf(splitted[0] == "nullable", true, false)
				values.Ignored = TerIf(splitted[0] == "ignored", true, false)
				continue
			}
			regexTagValueSplit := regexp.MustCompile(`(?sm)^(.*?):(.*?)$`)
			tagSplitted := regexTagValueSplit.FindStringSubmatch(item)
			if len(tagSplitted) > 2 {
				value := tagSplitted[2]
				values.In = TerIf(tagSplitted[1] == "in", value, values.In)
				values.Example = TerIf(tagSplitted[1] == "example", value, values.Example)
				values.Ref = TerIf(tagSplitted[1] == "$ref", value, values.Ref)
				values.MaxLength = TerIf(tagSplitted[1] == "maxLength", value, values.MaxLength)
				values.MinLength = TerIf(tagSplitted[1] == "minLength", value, values.MinLength)
				values.EnumValue = TerIf(tagSplitted[1] == "enumValue", value, values.EnumValue)
				values.Minimum = TerIf(tagSplitted[1] == "minimum", value, values.Minimum)
				values.Maximum = TerIf(tagSplitted[1] == "maximum", value, values.Maximum)
				values.Pattern = TerIf(tagSplitted[1] == "pattern", value, values.Pattern)
			}
		}
	}
	return values
}

func SanitizeCommentLineText(str string) string {
	str = RemoveNewLines(str)
	return str
}

type MergeMapType interface {
	Schema | Operations | Response
}

func MergeMaps[T MergeMapType](src map[string]T, dest map[string]T) map[string]T {
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

func MergePaths(src map[string]Operations, dest map[string]Operations) map[string]Operations {
	for key, value := range src {
		// merge different http methods for common path into one dict
		if _, ok := dest[key]; ok {
			dest[key] = Operations{
				Get:    TerIf(dest[key].Get != nil, dest[key].Get, value.Get),
				Put:    TerIf(dest[key].Put != nil, dest[key].Put, value.Put),
				Post:   TerIf(dest[key].Post != nil, dest[key].Post, value.Post),
				Delete: TerIf(dest[key].Delete != nil, dest[key].Delete, value.Delete),
				Patch:  TerIf(dest[key].Patch != nil, dest[key].Patch, value.Patch),
			}
		} else {
			dest[key] = value
		}
	}
	return dest
}

func GetResponseDescription(statusCode string) string {
	if desc, ok := ResponseDescriptions[statusCode]; ok {
		return desc
	}
	return "Unknown response code"
}

func GenerateOperationId(method string, path string) string {

	return ToUpperFirstLetter(strings.Split(path, "/")[1]) + RestActions[strings.ToUpper(method)] + RestOperations[strings.ToUpper(method)]
}
