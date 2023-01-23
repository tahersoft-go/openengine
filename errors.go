package openengine

import (
	"fmt"

	"gitlab.hoitek.fi/openapi/openengine/constants"
	"gitlab.hoitek.fi/openapi/openengine/engine"
)

func (p *openEngine) AddErrorResponses(errorResponses engine.ErrorResponses, errorResponseRefs ...string) OpenEngine {
	if len(p.Paths) != 0 {
		p.err = engine.BuildError("AddErrorResponses", "Paths already parsed. please add error responses before parsing paths")
		return p
	}

	var defaultRef string
	var normalizedErrorResponses = engine.ErrorResponses{}

	if len(errorResponseRefs) > 0 {
		defaultRef = engine.TerIf(errorResponseRefs[0] != "", "#/components/schemas/"+errorResponseRefs[0], "")
	}

	for statusCode, response := range errorResponses {
		if defaultRef != "" || response.Content.ApplicationJson.Schema.Ref != "" {
			response.Content.ApplicationJson.Schema.Ref = engine.TerIf(
				response.Content.ApplicationJson.Schema.Ref == "", defaultRef, "#/components/schemas/"+response.Content.ApplicationJson.Schema.Ref,
			)
		}
		if defaultRef != "" || response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref != "" {
			response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref = engine.TerIf(
				response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref == "", defaultRef, "#/components/schemas/"+response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref,
			)
		}
		normalizedErrorResponses[statusCode] = response
	}

	p.ErrorResponses = engine.MergeMaps(p.ErrorResponses, normalizedErrorResponses)
	return p
}

func (p *openEngine) AddDefaultErrors(codes ...int) OpenEngine {
	if len(p.Paths) != 0 {
		p.err = engine.BuildError("AddErrorResponses", "Paths already parsed. please add error responses before parsing paths")
		return p
	}

	if len(p.ErrorResponses) != 0 {
		p.err = engine.BuildError("AddErrorResponses", "Custom error responses already provided. please enable default error responses before adding your own ErrorResponsesDict")
		return p
	}

	// return filteredCodes
	var filteredCodes = engine.ErrorResponses{}
	for _, code := range codes {
		// change int to map
		if response, ok := constants.DefaultErrorResponses[fmt.Sprint(code)]; ok {
			filteredCodes[fmt.Sprint(code)] = response
		}
	}

	p.ErrorResponses = engine.MergeMaps(p.ErrorResponses, constants.DefaultErrorResponses)

	// if custom error responses are provided, filter default error responses
	if len(filteredCodes) != 0 {
		p.ErrorResponses = filteredCodes
	}

	return p
}
