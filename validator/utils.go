package validator

import (
	"errors"
	"regexp"

	"gopkg.in/yaml.v2"
)

func BuildError(errs *Errors, message string) {
	*errs = append(*errs, errors.New(message))
}

func HasSliceDuplicateString(slice []string) ([]string, bool) {
	hasDuplicatedValue := false
	duplicatedValues := []string{}
	for _, str := range slice {
		occur := 0
		for _, strVal := range slice {
			if str == strVal {
				occur++
				if occur > 1 {
					duplicatedValues = append(duplicatedValues, str)
				}
			}
			if occur > 1 {
				break
			}
		}
		if occur > 1 {
			hasDuplicatedValue = true
			break
		}
	}
	return duplicatedValues, hasDuplicatedValue
}

func FindAllByRegex(content, regex string) [][]string {
	regexCompiled := regexp.MustCompile(regex)
	return regexCompiled.FindAllStringSubmatch(content, -1)
}

func GetMapKeys(mapData map[string]interface{}) []string {
	keys := []string{}
	for key := range mapData {
		keys = append(keys, key)
	}
	return keys
}

func ToMapStringInterface(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
