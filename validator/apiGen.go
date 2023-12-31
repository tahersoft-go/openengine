package validator

import "github.com/tahersoft-go/openengine/engine"

type apiGen struct {
	Err       error
	RawResult string
	// Ignored Directories to search
	Ignore []string `yaml:"-"`
	// ErrorResponses
	ErrorResponses engine.ErrorResponses `yaml:"-"`
	// Escape yaml interpretation with -
	IsAutoTag bool `yaml:"-"`
	// OpenAPI data
	OpenApi      string              `yaml:"openapi"`
	Info         engine.Info         `yaml:"info"`
	ExternalDocs engine.ExternalDocs `yaml:"externalDocs,omitempty"`
	Servers      []engine.ApiServer  `yaml:"servers"`
	Tags         []engine.Tag        `yaml:"tags"`
	Paths        engine.PathsDict    `yaml:"paths"`
	Components   engine.Components   `yaml:"components"`
}
