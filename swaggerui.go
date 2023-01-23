package openengine

import (
	"log"
	"path"

	"gitlab.hoitek.fi/healthcare/services/maja/openengine/engine"
	"gitlab.hoitek.fi/healthcare/services/maja/openengine/swaggerui"
)

func (p *openEngine) ExportSwaggerUi(config engine.SwaggerUiConfig) OpenEngine {
	assetsPath := path.Join(config.ExportPath, "assets")
	log.Println(assetsPath)

	cssConfig := engine.AssetsConfig{
		ExportPath: assetsPath,
		FileName:   "swagger-ui.css",
		Link:       swaggerui.CssLink,
	}

	jsConfig := engine.AssetsConfig{
		ExportPath: assetsPath,
		FileName:   "swagger-ui-bundle.js",
		Link:       swaggerui.JsLink,
	}

	err := swaggerui.CreateFolderIfNotExists(assetsPath)
	if err != nil {
		p.err = err
		return p
	}

	_, err = swaggerui.WriteAsset(cssConfig)
	if err != nil {
		p.err = err
		return p
	}

	_, err = swaggerui.WriteAsset(jsConfig)
	if err != nil {
		p.err = err
		return p
	}

	htmlConfig := engine.HtmlConfig{
		ExportPath:      config.ExportPath,
		HtmlFileName:    config.HtmlFileName,
		Title:           config.Title,
		OpenApiFilePath: path.Join(config.ServeURI, p.fileName),
		CssPath:         path.Join(config.ServeURI, "assets", "swagger-ui.css"),
		JsPath:          path.Join(config.ServeURI, "assets", "swagger-ui-bundle.js"),
	}

	err = swaggerui.WriteHtml(htmlConfig)
	if err != nil {
		p.err = err
		return p
	}

	return p
}
