package openengine

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/tahersoft-go/openengine/engine"
)

var WgPath = &sync.WaitGroup{}

func (p *openEngine) extractPathsDataFromComments(handlersFilePath string) ([]engine.PathData, error) {
	// Create the AST by parsing src.
	fileSet := token.NewFileSet()

	// modelNames is a list of all the models in the file
	var pathsData []engine.PathData

	// Parse the file containing this struct definitions
	f, err := parser.ParseFile(fileSet, handlersFilePath, nil, parser.ParseComments)
	if err != nil {
		return pathsData, err
	}
	// If the file has no comments we return an error
	if len(f.Comments) == 0 {
		return pathsData, engine.BuildError(filepath.Base(handlersFilePath), "reading comment file: there is no @api declarations in the file")
	}
	// Loop through all the comments in the file
	for _, comment := range f.Comments {
		// log.Printf("comment %#v\n", comment)
		// If the comment list is empty we return an error
		if len(comment.List) == 0 {
			return pathsData, engine.BuildError(filepath.Base(handlersFilePath), "reading comment.List: there is no @api declarations in the file")
		}
		// Loop through all the comment lines from the list
		for _, commentLine := range comment.List {
			// log.Println("commentLine.Text", commentLine.Text)
			// remove new lies from text
			commentLineText := engine.SanitizeCommentLineText(commentLine.Text)
			// log.Println("commentLineText", commentLineText)
			// optional capture second item group with @apiDefine
			commentDataRegexp := regexp.MustCompile(engine.API_PATHS_DATA_REGEXP)
			// get second capture from regexp
			commentDataResult := commentDataRegexp.FindAllStringSubmatch(commentLineText, -1)
			// log.Printf("commentDataResult %#v\n", commentDataResult)
			pathData := engine.PathData{}
			if pathData.ApiSecurities == nil {
				pathData.ApiSecurities = map[string][]string{}
			}
			for _, comment := range commentDataResult {
				// FIXME: validate comment[1] based on other comment[2] value
				if len(comment) < 3 {
					return pathsData, errors.New("document comment is not valid: example @apiPath: /users")
				}
				comment[1] = strings.TrimSpace(comment[1])
				comment[2] = strings.TrimSpace(comment[2])
				// TODO: validate comment[2] based on other comment[1]
				switch comment[1] {
				case "@apiPath":
					pathData.ApiPath = comment[2]
				case "@apiMethod":
					pathData.ApiMethod = comment[2]
				case "@apiDescription":
					pathData.ApiDescription = comment[2]
				case "@apiSummary":
					pathData.ApiSummary = comment[2]
				case "@apiResponseRef":
					pathData.ApiResponseRef = comment[2]
				case "@apiRequestRef":
					pathData.ApiRequestRef = comment[2]
				case "@apiStatusCode":
					pathData.ApiStatusCode = comment[2]
				case "@apiTag":
					pathData.ApiTag = comment[2]
				case "@apiParametersRef":
					pathData.ApiParametersRef = comment[2]
				case "@apiDeprecated":
					pathData.ApiDeprecated = comment[2]
				case "@apiSecurity":
					scopesList := engine.TrimItemsSpace(strings.Split(comment[2], ","))
					securityName := scopesList[0]
					// escape true from 0 index and get the rest
					scopesList = scopesList[1:]
					pathData.ApiSecurities[securityName] = scopesList
				case "@apiErrorStatusCodes":
					pathData.ApiErrorStatusCodes =
						engine.TrimItemsSpace(strings.Split(comment[2], ","))
				}
				customRefRegexp := regexp.MustCompile(engine.API_CUSTOM_REF_REGEXP)
				customRefResult := customRefRegexp.FindStringSubmatch(comment[1])
				if len(customRefResult) == 2 {
					if pathData.ApiCustomErrorRefs == nil {
						pathData.ApiCustomErrorRefs = map[string]string{}
					}
					pathData.ApiCustomErrorRefs[customRefResult[1]] = comment[2]

				}

				customDescRegexp := regexp.MustCompile(engine.API_CUSTOM_DESCRIPTION_REGEXP)
				customDescResult := customDescRegexp.FindStringSubmatch(comment[1])
				if len(customDescResult) == 2 {
					if pathData.ApiCustomErrorDescriptions == nil {
						pathData.ApiCustomErrorDescriptions = map[string]string{}
					}
					pathData.ApiCustomErrorDescriptions[customDescResult[1]] = comment[2]

				}
			}
			// If we have a second capture we append it to the modelNames list
			pathsData = append(pathsData, pathData)
		}
	}

	// If we modelNames are empty we don't have any model, so we return error
	if len(pathsData) == 0 {
		return pathsData, engine.BuildError(filepath.Base(handlersFilePath), "parsing commentLine: there is no @api declarations in the file")
	}

	// return the modelNames list
	// log.Printf("pathsData %#v\n", pathsData)
	return pathsData, nil
}

func (p *openEngine) extractPathsDictFromFile(handlersFilePath string) (engine.PathsDict, error) {
	var pathsDict = engine.PathsDict{}

	commentsData, err := p.extractPathsDataFromComments(handlersFilePath)

	if err != nil {
		return nil, err
	}
	// log.Printf("commentsData: %#v\n", commentsData)

	if err != nil {
		return nil, err
	}

	for _, commentData := range commentsData {
		apiPath := engine.AddLeadingSlash(commentData.ApiPath)
		if apiPath == "" {
			continue
		}
		parameters := engine.Parameters{}
		parameterSchema, ok := p.Components.Schemas[commentData.ApiParametersRef]
		if ok {
			for name := range parameterSchema.Properties {
				parameters = append(parameters, engine.Parameter{
					// Description: "",
					Name:     name,
					In:       parameterSchema.Properties[name].In,
					Required: engine.TerIf(parameterSchema.Properties[name].In == "path", true, false),
					Example:  parameterSchema.Properties[name].Example,
					Schema: engine.ParameterSchema{
						Type: engine.TerIf(parameterSchema.Properties[name].Ref == "", parameterSchema.Properties[name].Type, ""),
						Ref:  engine.TerIf(parameterSchema.Properties[name].Ref != "", parameterSchema.Properties[name].Ref, ""),
					},
				})
			}
		}

		operation := engine.Operation{
			OperationId: engine.GenerateOperationId(commentData.ApiMethod, apiPath),
			Tags:        engine.TerIf(commentData.ApiTag != "", []string{commentData.ApiTag}, []string{}),
			Parameters:  parameters,
			Deprecated:  commentData.ApiDeprecated == "true",
			Description: commentData.ApiDescription,
			Summary:     commentData.ApiSummary,
			Responses:   engine.Responses{},
			Security:    engine.Security{},
		}
		if commentData.ApiRequestRef != "" {
			operation.RequestBody = &engine.RequestBody{
				Required: true,
				Content: engine.Content{
					ApplicationJson: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiRequestRef,
						},
					},
					ApplicationXWwwFormUrlencoded: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiRequestRef,
						},
					},
					MultipartFormData: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiRequestRef,
						},
					},
				},
			}
		}
		// Add Provided Default Errors from user

		if commentData.ApiResponseRef != "" {
			operation.Responses[commentData.ApiStatusCode] = engine.Response{
				Description: engine.GetResponseDescription(commentData.ApiStatusCode),
				Content: engine.Content{
					ApplicationJson: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiResponseRef,
						},
					},
					ApplicationXWwwFormUrlencoded: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiResponseRef,
						},
					},
					MultipartFormData: engine.MediaType{
						Schema: engine.DataSchema{
							Ref: "#/components/schemas/" + commentData.ApiResponseRef,
						},
					},
				},
			}
		}

		if len(p.ErrorResponses) > 0 {
			for statusCode, response := range p.ErrorResponses {
				codes := strings.Join(commentData.ApiErrorStatusCodes, ",")
				if !strings.Contains(codes, statusCode) && len(commentData.ApiErrorStatusCodes) > 0 {
					continue
				}
				operation.Responses[statusCode] = response
			}

			// handle customErrorRefs
			for statusCode, customErrorRef := range commentData.ApiCustomErrorRefs {
				// if we don't support statusCode we continue
				response, ok := operation.Responses[statusCode]
				if !ok {
					continue
				}
				ref := "#/components/schemas/" + customErrorRef
				customDescription, hasCustomDescription := commentData.ApiCustomErrorDescriptions[statusCode]
				description := engine.TerIf(hasCustomDescription, customDescription, engine.GetResponseDescription(statusCode))
				response.Description = description
				response.Content.ApplicationJson.Schema.Ref = ref
				response.Content.ApplicationXWwwFormUrlencoded.Schema.Ref = ref
				response.Content.MultipartFormData.Schema.Ref = ref
				operation.Responses[statusCode] = response
			}
		}

		for flow, scopes := range commentData.ApiSecurities {
			securityFlow := engine.SecurityFlow{}
			securityFlow[flow] = scopes
			operation.Security = append(operation.Security, securityFlow)
		}

		operation.Parameters = parameters
		if _, ok := pathsDict[apiPath]; !ok {
			pathsDict[apiPath] = engine.Operations{}
		}
		pathsDict[apiPath] = engine.Operations{
			Get:    engine.TerIf(strings.ToUpper(commentData.ApiMethod) == "GET", &operation, pathsDict[apiPath].Get),
			Put:    engine.TerIf(strings.ToUpper(commentData.ApiMethod) == "PUT", &operation, pathsDict[apiPath].Put),
			Post:   engine.TerIf(strings.ToUpper(commentData.ApiMethod) == "POST", &operation, pathsDict[apiPath].Post),
			Delete: engine.TerIf(strings.ToUpper(commentData.ApiMethod) == "DELETE", &operation, pathsDict[apiPath].Delete),
			Patch:  engine.TerIf(strings.ToUpper(commentData.ApiMethod) == "PATCH", &operation, pathsDict[apiPath].Patch),
		}
	}
	return pathsDict, nil
}

func (p *openEngine) extractPathsFromDirectory(handlersDirPath string, chanPaths chan engine.ChanPaths) (engine.PathsDict, error) {
	// Done waitgroup for this goroutine
	defer WgPath.Done()

	// Create AllPathsDict
	var AllPathsDict = engine.PathsDict{}

	// Get all the files in the models directory
	files, err := os.ReadDir(handlersDirPath)

	// If we have an error, we return it
	if err != nil {
		// Send ChanSchemas to channel if we have channel
		if chanPaths != nil {
			chanPaths <- engine.ChanPaths{
				Items: AllPathsDict,
				Err:   err,
			}
		}

		// Return AllSchemasDict
		return AllPathsDict, err
	}

	// Loop through all the files
	for _, file := range files {

		// If the file is a directory or is an ignored file we continue
		if file.IsDir() || engine.IsIgnoredFile(file.Name()) {
			continue
		}

		// Extract the schemas from the file
		pathsDict, err := p.extractPathsDictFromFile(handlersDirPath + "/" + file.Name())

		// If we have an error, we return it
		if err != nil {
			// TODO: we should have verbose param to log more
			// log.Println(err)
			continue
		}

		// Append the schemas to the global schemas map
		AllPathsDict = engine.MergeMaps(AllPathsDict, pathsDict)
	}
	// Send ChanSchemas to channel if we have channel
	if chanPaths != nil {
		chanPaths <- engine.ChanPaths{
			Items: AllPathsDict,
			Err:   nil,
		}
	}

	// Return AllSchemasDict
	return AllPathsDict, nil
}

func (p *openEngine) AddPaths(pathsDict engine.PathsDict) OpenEngine {
	p.Paths = engine.MergeMaps(p.Paths, pathsDict)
	return p
}

func (p *openEngine) ParsePaths(baseDirectory string, allIgnoredPaths ...[]string) OpenEngine {
	if p.Components.Schemas == nil || len(p.Components.Schemas) == 0 {
		p.err = errors.New("parse your schemas first")
		return p
	}

	// Create Mutex for SchemasDict
	mx := &sync.Mutex{}

	// Initialize WaitGroup for goroutines
	WgPath = &sync.WaitGroup{}

	// Create SchemasDict
	AllPathsDict := engine.PathsDict{}

	// Find all the directories in the baseDirPath
	var ignoredPaths []string
	if len(allIgnoredPaths) > 0 {
		ignoredPaths = allIgnoredPaths[0]
	}
	mergedIgnoredPaths := append(p.GeneralIgnoredPaths, ignoredPaths...)
	handlersDirectoryPaths, err := engine.FindAllDirectoriesInPath(baseDirectory, &mergedIgnoredPaths)

	// If we have an error, we return it
	if err != nil {
		p.err = err
		return p
	}

	// Create a channel for the schemas
	var chanPaths = make(chan engine.ChanPaths, len(handlersDirectoryPaths))

	// Close the channel when we are done with it in goroutine
	go func() {
		defer close(chanPaths)
		WgPath.Wait()
	}()

	// Loop through all the models paths
	for _, handlersDirectoryPath := range handlersDirectoryPaths {
		// Add 1 to WaitGroup
		WgPath.Add(1)

		// Extract the schemas from the directory in a goroutine
		go p.extractPathsFromDirectory(handlersDirectoryPath, chanPaths)
	}

	// Loop through all the schemas in the channel
	for result := range chanPaths {
		// Lock the mutex for the SchemasDict
		mx.Lock()

		// Check if we have an error
		if result.Err != nil {
			p.err = result.Err
			return p
		}

		// Merge the schemas to the global schemas map
		AllPathsDict = engine.MergeMaps(AllPathsDict, result.Items)

		// Unlock the mutex for the SchemasDict
		mx.Unlock()
	}

	// Return the global schemas map

	p.Paths = AllPathsDict
	return p
}
