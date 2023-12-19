package openengine

import (
	"go/ast"
	"path"

	"github.com/tahersoft-go/openengine/engine"
	"github.com/tahersoft-go/openengine/validator"
	"gopkg.in/yaml.v2"
)

type Init struct {
	Info         engine.Info         `yaml:"info"`
	ExternalDocs engine.ExternalDocs `yaml:"externalDocs,omitempty"`
}

type openEngine struct {
	// Escape yaml interpretation with -
	// Error
	err error
	// Yaml Result string
	rawResult string
	// FileName
	fileName string
	// GeneralIgnoredPaths Directories to search
	GeneralIgnoredPaths []string `yaml:"-"`
	// Ignored directory names to search
	ErrorResponses engine.ErrorResponses `yaml:"-"`
	// OpenAPI data
	OpenApi      string              `yaml:"openapi"`
	Info         engine.Info         `yaml:"info"`
	ExternalDocs engine.ExternalDocs `yaml:"externalDocs,omitempty"`
	Servers      engine.ApiServers   `yaml:"servers"`
	Tags         []engine.Tag        `yaml:"tags"`
	Paths        engine.PathsDict    `yaml:"paths"`
	Components   engine.Components   `yaml:"components"`
}

type OpenEngine interface {
	// FileName
	SetFileName(fileName string) OpenEngine
	// Ignores
	AddIgnoredPaths(dirs []string) OpenEngine
	// Error Responses
	AddErrorResponses(errorResponses engine.ErrorResponses, defaultRef ...string) OpenEngine
	AddDefaultErrors(...int) OpenEngine
	// Tags
	AddTag(tag engine.Tag) OpenEngine
	ParseTags(domainsPath string, ignoredDirs ...[]string) OpenEngine
	// Servers
	AddServers(servers engine.ApiServers) OpenEngine
	//Schemas
	extractSchemaNamesFromComments(schemasFilePath string) ([]string, error)
	mapSchemaFieldsToSchemaDict(list []*ast.Field, structName string, schemasDict *engine.SchemasDict)
	extractSchemasDictFromFile(schemasFilePath string) (engine.SchemasDict, error)
	extractSchemasFromDirectory(structsDirPath string, chanSchemas chan engine.ChanSchemas) (engine.SchemasDict, error)
	AddSchemas(schemasDict engine.SchemasDict) OpenEngine
	ParseSchemas(path string, ignoredPaths ...[]string) OpenEngine
	//enums
	extractEnumNamesFromComments(schemasFilePath string) ([]string, error)
	mapEnumFieldsToSchemaDict(list []*ast.Field, structName string, schemasDict *engine.SchemasDict)
	extractEnumsDictFromFile(schemasFilePath string) (engine.SchemasDict, error)
	extractEnumsFromDirectory(structsDirPath string, chanSchemas chan engine.ChanSchemas) (engine.SchemasDict, error)
	AddEnums(schemasDict engine.SchemasDict) OpenEngine
	ParseEnums(path string, ignoredPaths ...[]string) OpenEngine
	// Paths
	extractPathsDataFromComments(handlersFilePath string) ([]engine.PathData, error)
	extractPathsDictFromFile(handlersFilePath string) (engine.PathsDict, error)
	extractPathsFromDirectory(handlersDirPath string, chanPaths chan engine.ChanPaths) (engine.PathsDict, error)
	AddPaths(pathsDict engine.PathsDict) OpenEngine
	ParsePaths(handlersDirsPaths string, ignoredPaths ...[]string) OpenEngine
	// SecuritySchemas
	AddSecuritySchemes(securitySchemas engine.SecuritySchemesTypes) OpenEngine
	// SwaggerUI
	ExportSwaggerUi(config engine.SwaggerUiConfig) OpenEngine
	// Final
	Generate(dest ...string) (string, error)
}

func NewPackage(data ...Init) OpenEngine {
	var init Init
	if len(data) > 0 {
		init = data[0]
	}

	info := engine.Info{
		Title:          engine.TerIf(init.Info.Title == "", engine.TITLE, init.Info.Title),
		Description:    engine.TerIf(init.Info.Description == "", engine.DESCRIPTION, init.Info.Description),
		Version:        engine.TerIf(init.Info.Version == "", engine.VERSION, init.Info.Version),
		TermsOfService: engine.TerIf(init.Info.TermsOfService == "", engine.TERMS_OF_SERVICE, init.Info.TermsOfService),
		Contact: engine.Contact{
			Name:  engine.TerIf(init.Info.Contact.Name == "", engine.CONTACT_NAME, init.Info.Contact.Name),
			Email: engine.TerIf(init.Info.Contact.Email == "", engine.CONTACT_EMAIL, init.Info.Contact.Email),
		},
		License: engine.License{
			Name: engine.TerIf(init.Info.License.Name == "", engine.LICENSE_NAME, init.Info.License.Name),
			Url:  engine.TerIf(init.Info.License.Url == "", engine.LICENSE_URL, init.Info.License.Url),
		},
	}

	externalDocs := engine.ExternalDocs{
		Description: engine.TerIf(init.ExternalDocs.Description == "", engine.EXTERNAL_DOCS_DESCRIPTION, init.ExternalDocs.Description),
		Url:         engine.TerIf(init.ExternalDocs.Url == "", engine.EXTERNAL_DOCS_URL, init.ExternalDocs.Url),
	}

	return &openEngine{
		fileName:            engine.DEFAULT_FILE_NAME,
		GeneralIgnoredPaths: engine.IgnoredDirectories,
		OpenApi:             engine.OPEN_API_VERSION,
		Info:                info,
		ExternalDocs:        externalDocs,
		Servers:             engine.ApiServers{},
		Tags:                []engine.Tag{},
		Paths:               engine.PathsDict{},
		Components: engine.Components{
			Schemas:         engine.SchemasDict{},
			RequestBodies:   engine.RequestBodies{},
			SecuritySchemes: engine.SecuritySchemes{},
		},
	}
}

func (p *openEngine) SetFileName(fileName string) OpenEngine {
	if fileName != "" {
		p.fileName = fileName
	}
	return p
}

func (p *openEngine) AddIgnoredPaths(dirs []string) OpenEngine {
	p.GeneralIgnoredPaths = append(p.GeneralIgnoredPaths, dirs...)
	return p
}

func (p *openEngine) Generate(destinationDirectories ...string) (string, error) {
	providedPath := p.fileName
	if len(destinationDirectories) > 0 {
		providedPath = path.Join(destinationDirectories[0], p.fileName)
	}

	if p.err != nil {
		return p.rawResult, p.err
	}

	yamlDocs, err := yaml.Marshal(p)
	if err != nil {
		p.err = err
		return p.rawResult, p.err
	}

	er := validator.ValidateRaw(string(yamlDocs))

	forceValidate := false

	if er != nil {
		p.err = err
		if forceValidate {
			return p.rawResult, p.err
		}
	}

	err = engine.ExportAPIDocsYaml(providedPath, string(yamlDocs))
	if err != nil {
		p.err = err
		return p.rawResult, p.err
	}
	p.rawResult = string(yamlDocs)

	return p.rawResult, p.err
}
