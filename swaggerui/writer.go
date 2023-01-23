package swaggerui

import (
	"io"
	"net/http"
	"os"
	"path"

	"gitlab.hoitek.fi/openapi/openengine/engine"
)

func CreateFolderIfNotExists(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func WriteHtml(config engine.HtmlConfig) error {
	fileName := engine.TerIf(config.HtmlFileName == "", "index.html", config.HtmlFileName)
	filePath := path.Join(config.ExportPath, fileName)

	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(HtmlTemplate(config))
	if err != nil {
		return err
	}

	return err
}

func WriteAsset(config engine.AssetsConfig) (string, error) {
	filePath := path.Join(config.ExportPath, config.FileName)

	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return filePath, nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	response, err := http.Get(config.Link)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, err
}
