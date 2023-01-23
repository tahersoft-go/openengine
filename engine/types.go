package engine

// ---------------SecuritySchemaTypes----------------
type ApiKeys string
type AuthType string

type ApiKeySecurityScheme struct {
	Type        AuthType `yaml:"type,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Name        string   `yaml:"name,omitempty"`
	In          ApiKeys  `yaml:"in,omitempty"`
}

type HttpSecuritySchemeType string
type HttpSecurityBearerFormat string

type HttpSecurityScheme struct {
	Type         AuthType                 `yaml:"type,omitempty"`
	Description  string                   `yaml:"description,omitempty"`
	Scheme       HttpSecuritySchemeType   `yaml:"scheme,omitempty"`
	BearerFormat HttpSecurityBearerFormat `yaml:"bearerFormat,omitempty"`
}

type OAuth2SecurityScheme struct {
	Type        AuthType    `yaml:"type,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Flows       Oauth2Flows `yaml:"flows,omitempty"`
}

type Oauth2Flows struct {
	// Implicit Flow without PKCE and refresh token because no redirect uri is supported for old browsers
	Implicit *Oauth2ImplicitFlow `yaml:"implicit,omitempty"`
	// High Trusted Clients like ssr or mobile apps
	ResourceOwnerPassword *Oauth2PasswordFlow `yaml:"password,omitempty"`
	// Machine to Machine
	ClientCredentials *Oauth2CredsFlow `yaml:"clientCredentials,omitempty"`
	// Low Trusted Clients with external popup resource owner data handling
	AuthorizationCodeWithPKCE *Oauth2CodeFlow `yaml:"authorizationCode,omitempty"`
	// Low Trusted Clients with internal resource owner data handling
	// InteractionCodeWithPKCE *Oauth2CodeFlow `yaml:"authorizationCode,omitempty"`
}

type Oauth2ImplicitFlow struct {
	AuthorizationUrl string       `yaml:"authorizationUrl,omitempty"`
	RefreshUrl       string       `yaml:"refreshUrl,omitempty"`
	Scopes           OAuth2Scopes `yaml:"scopes,omitempty"`
}

type Oauth2PasswordFlow struct {
	TokenUrl   string       `yaml:"tokenUrl,omitempty"`
	RefreshUrl string       `yaml:"refreshUrl,omitempty"`
	Scopes     OAuth2Scopes `yaml:"scopes,omitempty"`
}

type Oauth2CredsFlow struct {
	TokenUrl   string       `yaml:"tokenUrl,omitempty"`
	RefreshUrl string       `yaml:"refreshUrl,omitempty"`
	Scopes     OAuth2Scopes `yaml:"scopes,omitempty"`
}

type Oauth2CodeFlow struct {
	AuthorizationUrl string       `yaml:"authorizationUrl,omitempty"`
	TokenUrl         string       `yaml:"tokenUrl,omitempty"`
	RefreshUrl       string       `yaml:"refreshUrl,omitempty"`
	Scopes           OAuth2Scopes `yaml:"scopes,omitempty"`
}

type OAuth2Scopes map[string]string

type OpenIdSecurityScheme struct {
	Type             AuthType `yaml:"type,omitempty"`
	Description      string   `yaml:"description,omitempty"`
	OpenIdConnectUrl string   `yaml:"openIdConnectUrl,omitempty"`
}

type SecuritySchemes map[string]interface{}
type SecuritySchemesTypes struct {
	ApiKey ApiKeySecuritySchemesDict
	Http   HttpSecuritySchemesDict
	OAuth2 OAuth2SecuritySchemesDict
	OpenId OpenIdSecuritySchemesDict
}

type ApiKeySecuritySchemesDict map[string]ApiKeySecurityScheme
type HttpSecuritySchemesDict map[string]HttpSecurityScheme
type OAuth2SecuritySchemesDict map[string]OAuth2SecurityScheme
type OpenIdSecuritySchemesDict map[string]OpenIdSecurityScheme

// -------------------------------------------------

type SecurityScopesList []string
type SecurityFlow map[string]SecurityScopesList

type (
	SchemasDict   map[string]Schema
	PathsDict     map[string]Operations
	Properties    map[string]Property
	RequestBodies map[string]RequestBody
	Responses     map[string]Response
	Security      []SecurityFlow
	Parameters    []Parameter
	Tags          []Tag
)
type PathData struct {
	ApiPath                    string
	ApiMethod                  string
	ApiDescription             string
	ApiSummary                 string
	ApiRequestRef              string
	ApiResponseRef             string
	ApiStatusCode              string
	ApiTag                     string
	ApiParametersRef           string
	ApiDeprecated              string
	ApiErrorStatusCodes        []string
	ApiCustomErrorRefs         map[string]string
	ApiCustomErrorDescriptions map[string]string
	ApiSecurities              map[string][]string
}

type OpenApiFieldTagValues struct {
	In        string `yaml:"in,omitempty"`
	Example   string `yaml:"example,omitempty"`
	Ref       string `yaml:"$ref,omitempty"`
	Required  bool   `yaml:"required,omitempty"`
	Nullable  bool   `yaml:"nullable,omitempty"`
	MaxLength string `yaml:"maxLength,omitempty"`
	MinLength string `yaml:"minLength,omitempty"`
	Minimum   string `yaml:"minimum,omitempty"`
	Maximum   string `yaml:"maximum,omitempty"`
	Pattern   string `yaml:"pattern,omitempty"`
	Ignored   bool   `yaml:"ignored,omitempty"`
}

type ErrorResponses Responses

type ChanSchemas struct {
	Items SchemasDict
	Err   error
}

type ChanPaths struct {
	Items PathsDict
	Err   error
}

type Contact struct {
	Name  string `yaml:"name,omitempty"`
	Url   string `yaml:"url,omitempty"`
	Email string `yaml:"email,omitempty"`
}

type PropertyItems struct {
	Ref string `yaml:"$ref,omitempty"`
}

type Property struct {
	In        string         `yaml:"-"`
	Type      string         `yaml:"type,omitempty"`
	Format    string         `yaml:"format,omitempty"`
	Example   string         `yaml:"example,omitempty"`
	Ref       string         `yaml:"$ref,omitempty"`
	Pattern   string         `yaml:"pattern,omitempty"`
	Items     *PropertyItems `yaml:"items,omitempty"`
	Required  bool           `yaml:"required,omitempty"`
	Nullable  bool           `yaml:"nullable,omitempty"`
	MaxLength int            `yaml:"maxLength,omitempty"`
	MinLength int            `yaml:"minLength,omitempty"`
	Minimum   int            `yaml:"minimum,omitempty"`
	Maximum   int            `yaml:"maximum,omitempty"`
}

type License struct {
	Name string `yaml:"name,omitempty"`
	Url  string `yaml:"url,omitempty"`
}

type ExternalDocs struct {
	Description string `yaml:"description,omitempty"`
	Url         string `yaml:"url,omitempty"`
}

type ApiServers []ApiServer

type ApiServer struct {
	Url string `yaml:"url,omitempty"`
}

type Tag struct {
	Name string `yaml:"name,omitempty"`
}

type Info struct {
	Title          string  `yaml:"title,omitempty"`
	Description    string  `yaml:"description,omitempty"`
	Version        string  `yaml:"version,omitempty"`
	TermsOfService string  `yaml:"termsOfService,omitempty"`
	Contact        Contact `yaml:"contact,omitempty"`
	License        License `yaml:"license,omitempty"`
}

type DataSchema struct {
	Ref string `yaml:"$ref,omitempty"`
}

type ParameterSchema struct {
	Type    string   `yaml:"type,omitempty"`
	Default string   `yaml:"default,omitempty"`
	Enum    []string `yaml:"enum,omitempty"`
	Ref     string   `yaml:"$ref,omitempty"`
}

type MediaType struct {
	Schema DataSchema `yaml:"schema,omitempty"`
}

type Content struct {
	ApplicationJson               MediaType `yaml:"application/json,omitempty"`
	ApplicationXWwwFormUrlencoded MediaType `yaml:"application/x-www-form-urlencoded,omitempty"`
}

type Schema struct {
	Type       string     `yaml:"type,omitempty"`
	Format     string     `yaml:"format,omitempty"`
	Properties Properties `yaml:"properties,omitempty"`
	Required   []string   `yaml:"required,omitempty"`
}

// Response
type Response struct {
	Description string  `yaml:"description,omitempty"`
	Content     Content `yaml:"content,omitempty"`
}

// Request
type RequestBody struct {
	Description string  `yaml:"description,omitempty"`
	Content     Content `yaml:"content,omitempty"`
	Required    bool    `yaml:"required,omitempty"`
}

type Components struct {
	Schemas         SchemasDict            `yaml:"schemas,omitempty"`
	RequestBodies   RequestBodies          `yaml:"requestBodies,omitempty"`
	SecuritySchemes map[string]interface{} `yaml:"securitySchemes,omitempty"`
}

type Parameter struct {
	Name        string          `yaml:"name,omitempty"`
	In          string          `yaml:"in,omitempty"`
	Description string          `yaml:"description,omitempty"`
	Required    bool            `yaml:"required,omitempty"`
	Schema      ParameterSchema `yaml:"schema,omitempty"`
	Example     string          `yaml:"example,omitempty"`
}

type Operation struct {
	Tags        []string     `yaml:"tags,omitempty"`
	Summary     string       `yaml:"summary,omitempty"`
	Description string       `yaml:"description,omitempty"`
	OperationId string       `yaml:"operationId,omitempty"`
	Parameters  Parameters   `yaml:"parameters,omitempty"`
	RequestBody *RequestBody `yaml:"requestBody,omitempty"`
	Responses   Responses    `yaml:"responses,omitempty"`
	Security    Security     `yaml:"security,omitempty"`
	Deprecated  bool         `yaml:"deprecated,omitempty"`
}

type Operations struct {
	Put    *Operation `yaml:"put,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Get    *Operation `yaml:"get,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty"`
}

type SwaggerUiConfig struct {
	Title        string
	ExportPath   string
	ServeURI     string
	HtmlFileName string
}

type HtmlConfig struct {
	HtmlFileName    string
	ExportPath      string
	Title           string
	CssPath         string
	JsPath          string
	OpenApiFilePath string
}

type AssetsConfig struct {
	ExportPath string
	FileName   string
	Link       string
}
