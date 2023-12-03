package validator

import (
	"fmt"
	"strings"

	"github.com/tahersoft-go/openengine/engine"
)

func (v *openApiValidator) checkResponseRefExistsInSchema(responses engine.Responses) {
	for _, response := range responses {
		ref := response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref
		if ref != "" {
			if !v.isRefExistsInSchema(ref) {
				BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
			}
		}
		ref = response.Content.ApplicationJson.Schema.Ref
		if ref != "" {
			if !v.isRefExistsInSchema(ref) {
				BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
			}
		}
	}
}

func (v *openApiValidator) isRefExistsInSchema(ref string) bool {
	splittedRefName := strings.Split(ref, "/")
	refName := splittedRefName[len(splittedRefName)-1]
	for schemaName, _ := range v.YamlDoc.Components.Schemas {
		if schemaName == refName {
			return true
		}
	}
	return false
}

func (v *openApiValidator) CheckAllRefsExistsInSchema() *openApiValidator {
	// Check path refs
	for _, operations := range v.YamlDoc.Paths {
		if operations.Get != nil {
			if operations.Get.RequestBody != nil {
				ref := operations.Get.RequestBody.Content.ApplicationJson.Schema.Ref
				if ref != "" {
					if !v.isRefExistsInSchema(ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
					}
				}
			}
			v.checkResponseRefExistsInSchema(operations.Get.Responses)
		}
		if operations.Put != nil {
			if operations.Put.RequestBody != nil {
				ref := operations.Put.RequestBody.Content.ApplicationJson.Schema.Ref
				if ref != "" {
					if !v.isRefExistsInSchema(ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
					}
				}
			}
			v.checkResponseRefExistsInSchema(operations.Put.Responses)
		}
		if operations.Patch != nil {
			if operations.Patch.RequestBody != nil {
				ref := operations.Patch.RequestBody.Content.ApplicationJson.Schema.Ref
				if ref != "" {
					if !v.isRefExistsInSchema(ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
					}
				}
			}
			v.checkResponseRefExistsInSchema(operations.Patch.Responses)
		}
		if operations.Delete != nil {
			if operations.Delete.RequestBody != nil {
				ref := operations.Delete.RequestBody.Content.ApplicationJson.Schema.Ref
				if ref != "" {
					if !v.isRefExistsInSchema(ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
					}
				}
			}
			v.checkResponseRefExistsInSchema(operations.Delete.Responses)
		}
		if operations.Post != nil {
			if operations.Post.RequestBody != nil {
				ref := operations.Post.RequestBody.Content.ApplicationJson.Schema.Ref
				if ref != "" {
					if !v.isRefExistsInSchema(ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", ref))
					}
				}
			}
			v.checkResponseRefExistsInSchema(operations.Post.Responses)
		}
	}

	// Check Schema refs
	for _, schema := range v.YamlDoc.Components.Schemas {
		if schema.Properties != nil {
			for _, prop := range schema.Properties {
				if prop.Ref != "" {
					if !v.isRefExistsInSchema(prop.Ref) {
						BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", prop.Ref))
					}
				}
				if prop.Items != nil {
					if prop.Items.Ref != "" {
						if !v.isRefExistsInSchema(prop.Items.Ref) {
							BuildError(&v.Errors, fmt.Sprintf("Ref %s does not exist in schema", prop.Items.Ref))
						}
					}
				}
			}
		}
	}
	return v
}
