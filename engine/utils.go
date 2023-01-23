package engine

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func ExtractDirNames(internalFolderPath string) ([]string, error) {
	folders, err := os.ReadDir(internalFolderPath)
	if err != nil {
		return []string{}, err
	}

	var tags []string

	for _, folder := range folders {
		if folder.IsDir() {
			tags = append(tags, folder.Name())
		}
	}

	return tags, nil
}

func TerIf[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func TerIfNil[T any](cond, vtrue *T) *T {
	if cond != nil {
		return cond
	}
	return vtrue
}

// recurse find all dirs in path
func FindAllDirectoriesInPath(path string, ignoredDirs *[]string) ([]string, error) {
	var dirs []string
	files, err := os.ReadDir(path)

	if err != nil {
		return dirs, err
	}
	for _, file := range files {
		fileName := file.Name()
		if StringInSlice(fileName, ignoredDirs) {
			continue
		}
		if file.IsDir() {
			dirs = append(dirs, path+"/"+fileName)
			subDirs, err := FindAllDirectoriesInPath(path+"/"+fileName, ignoredDirs)
			if err != nil {
				log.Println(err)
				continue
			}
			dirs = append(dirs, subDirs...)
		}
	}
	return dirs, nil
}

func StringInSlice(a string, list *[]string) bool {
	for _, b := range *list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveNewLines(str string) string {
	return strings.ReplaceAll(str, "\n", " ")
}

func RemoveSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func ToUpperFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func BuildLogText(label, description string) string {
	return fmt.Sprintf("\n------\n%s: %s\n------\n", label, description)
}

func BuildLog(label, description string) {
	log.Println(BuildLogText(label, description))
}

func BuildError(label, description string) error {
	return errors.New(BuildLogText(label, description))
}

func TrimItemsSpace(items []string) []string {
	var result []string
	for _, item := range items {
		result = append(result, strings.TrimSpace(item))
	}
	return result
}

func AddLeadingSlash(path string) string {
	if path == "" {
		return ""
	}

	if path[0] != '/' {
		return "/" + path
	}
	return path
}
