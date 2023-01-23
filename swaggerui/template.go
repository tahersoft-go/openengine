package swaggerui

import "gitlab.hoitek.fi/openapi/openengine/engine"

func HtmlTemplate(config engine.HtmlConfig) string {
	return `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta
      name="` + config.Title + `"
      content="SwaggerUI"
    />
    <title>` + config.Title + `</title>
    <link rel="stylesheet" href="` + config.CssPath + `" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="` + config.JsPath + `" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '` + config.OpenApiFilePath + `',
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
        ],
      });
    };
  </script>
  </body>
</html>
	`
}
