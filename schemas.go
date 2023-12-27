package openengine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/tahersoft-go/openengine/engine"
)

var WgSchema = &sync.WaitGroup{}

func (p *openEngine) extractSchemaNamesFromComments(schemasFilePath string) ([]string, error) {
	// Create the AST by parsing src.
	fileSet := token.NewFileSet()

	// structNames is a list of all the models in the file
	var structNames []string

	// Parse the file containing this struct definitions
	f, err := parser.ParseFile(fileSet, schemasFilePath, nil, parser.ParseComments)
	if err != nil {
		return structNames, err
	}
	fileName := filepath.Base(schemasFilePath)
	// If the file has no comments we return an error
	if len(f.Comments) == 0 {
		return structNames, engine.BuildError(fileName, "reading comment file: there is no comments (@api declarations) in the file")
	}
	// Loop through all the comments in the file
	for _, comment := range f.Comments {
		// log.Printf("comment %#v\n", comment)
		// If the comment list is empty we return an error
		if len(comment.List) == 0 {
			return structNames, engine.BuildError(fileName, "reading comment.List: there is no @api declarations in the file")
		}
		// log.Println("commentList", comment.List)
		// Loop through all the comment lines from the list
		for _, commentLine := range comment.List {
			// log.Println("commentLine", commentLine)
			// optional capture second item group with @apiDefine
			reg := regexp.MustCompile(engine.API_SCHEMAS_DATA_REGEXP)
			// Sanitize the comment line text and remove new lines and spaces to make regexp work
			commentLineText := engine.SanitizeCommentLineText(commentLine.Text)
			// get second capture from regexp
			structName := reg.FindStringSubmatch(commentLineText)
			// If we have a second capture we append it to the modelNames list
			if structName != nil {
				// If the second capture is not empty we append it to the modelNames list
				if len(structName) > 1 {
					structNames = append(structNames, structName[2])
					break
				}
			}
		}
	}

	// If we modelNames are empty we don't have any model, so we return error
	if len(structNames) == 0 {
		return structNames, engine.BuildError(fileName, "parsing commentLine: there is no @api declarations in the file")
	}

	// return the modelNames list
	// log.Println("structNames", structNames)
	return structNames, nil
}

// TODO: refactor this function
func (p *openEngine) mapSchemaFieldsToSchemaDict(list []*ast.Field, structName string, schemasDict *engine.SchemasDict) {
	for _, field := range list {
		// log.Printf("-----%s Struct -> %s %#v\n", structName, field.Names[0].Name, field.Type)
		if field.Tag == nil {
			continue
		}
		// Get the tag value and remove start and end quotes
		var fieldTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

		// openapi tag parsed
		openApiTagValues := fieldTag.Get(engine.OPEN_API_TAG_NAME)
		// openapi tag values in struct
		tagValues := engine.ParseStructTagValues(openApiTagValues)
		// if tag is ignored we continue
		if tagValues.Ignored {
			continue
		}
		// json tag parsed
		jsonFieldName := strings.Split(fieldTag.Get(engine.JSON_TAG_NAME), ",")[0]
		// json fieldName if json tag is not empty
		fieldName := engine.TerIf(jsonFieldName == "", field.Names[0].Name, strings.Split(jsonFieldName, ",")[0])

		var (
			tp     string
			format string
			ref    string

			items *engine.PropertyItems

			maxLength, _ = strconv.Atoi(tagValues.MaxLength)
			minLength, _ = strconv.Atoi(tagValues.MinLength)
			maximum, _   = strconv.Atoi(tagValues.Maximum)
			minimum, _   = strconv.Atoi(tagValues.Minimum)

			in = engine.TerIf(tagValues.In != "", tagValues.In, "query")
		)

		switch field.Type.(type) {
		case *ast.Ident:
			// Get the type of the field
			t := field.Type.(*ast.Ident).Name

			tp = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, engine.OpenAPITypes(t), "object")
			format = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, engine.OpenAPIFormats(t), "object")
			ref = engine.TerIf(field.Type.(*ast.Ident).Obj != nil, "#/components/schemas/"+tagValues.Ref, "")
			maxLength = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, maxLength, 0)
			minLength = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, minLength, 0)
			maximum = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, maximum, 0)
			minimum = engine.TerIf(field.Type.(*ast.Ident).Obj == nil, minimum, 0)

			// Set current property with specific modelName on SchemasDict

		case *ast.IndexExpr, *ast.SelectorExpr, *ast.ArrayType:
			if tagValues.Ref == "" {
				continue
			}

			tp = "object"
			format = "object"
			ref = "#/components/schemas/" + tagValues.Ref

			// if field is array type we set the items ref and clear outer ref
			if fmt.Sprintf("%T", field.Type) == "*ast.ArrayType" {
				tp = "array"
				format = "array"
				items = &engine.PropertyItems{
					Ref: ref,
				}
				ref = ""
			}
		default:
			log.Printf("field.Type: %#v is not supported yet", field.Type)
		}
		// Set current property with specific modelName on SchemasDict
		(*schemasDict)[structName].Properties[fieldName] = engine.Property{
			In:        in,
			Type:      engine.TerIf(ref == "", tp, ""),
			Format:    engine.TerIf(ref == "", format, ""),
			Example:   engine.TerIf(ref == "", tagValues.Example, ""),
			Nullable:  engine.TerIf(ref == "", tagValues.Nullable, false),
			Pattern:   engine.TerIf(ref == "", tagValues.Pattern, ""),
			MaxLength: engine.TerIf(ref == "", maxLength, 0),
			MinLength: engine.TerIf(ref == "", minLength, 0),
			Maximum:   engine.TerIf(ref == "", maximum, 0),
			Minimum:   engine.TerIf(ref == "", minimum, 0),
			Ref:       ref,
			Items:     items,
		}
	}

	// add required to fields with tag `required`
	for _, field := range list {
		if field.Tag != nil {
			// Get the tag value
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
			// Get the field name
			fieldName := field.Names[0].Name
			// If the field has a json tag, use that as the field name
			jsonFieldName := tag.Get(engine.JSON_TAG_NAME)
			fieldName = engine.TerIf(jsonFieldName == "", fieldName, strings.Split(jsonFieldName, ",")[0])
			// Get the example value
			openApiTagValues := tag.Get(engine.OPEN_API_TAG_NAME)
			tagValues := engine.ParseStructTagValues(openApiTagValues)
			schema := (*schemasDict)[structName]
			if tagValues.Required {
				schema.Required = append(schema.Required, fieldName)
			}
			(*schemasDict)[structName] = schema
		}
	}
}

func (p *openEngine) extractSchemasDictFromFile(schemasFilePath string) (engine.SchemasDict, error) {
	var schemasDict = engine.SchemasDict{}
	// Create the AST by parsing src.
	fileSet := token.NewFileSet()

	// Parse the file containing this struct definitions
	f, err := parser.ParseFile(fileSet, schemasFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// structNames, err := ExtractDataFromComments(schemasFilePath)
	structNames, err := p.extractSchemaNamesFromComments(schemasFilePath)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Loop through all the declarations in the file
	for _, decl := range f.Decls {
		// If the declaration is a GenDecl (general declaration)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Loop through all the Specs and grab the one we want.
			for _, spec := range genDecl.Specs {
				// If the spec is a TypeSpec we are half way there.
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// If the Type of the TypeSpec is a StructType, we found the struct
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// check if modelNames have the typeSpec.Name.Name, If not we continue
						if !engine.StringInSlice(typeSpec.Name.Name, &structNames) {
							continue
						}
						// Get struct name
						structName := typeSpec.Name.Name
						// Create a new Schema for the struct
						schemasDict[structName] = engine.Schema{
							// This is a hack to get the type of the struct
							Type:       "object",
							Format:     "object",
							Properties: engine.Properties{},
						}
						// Loop through all the fields in the struct
						p.mapSchemaFieldsToSchemaDict(structType.Fields.List, structName, &schemasDict)
					}
				}
			}
		}
	}
	// log.Println("SchemasDict", schemasDict)
	return schemasDict, nil
}

func (p *openEngine) extractSchemasFromDirectory(structsDirPath string, chanSchemas chan engine.ChanSchemas) (engine.SchemasDict, error) {
	// Done waitgroup for this goroutine
	defer WgSchema.Done()

	// Create AllSchemasDict
	var AllSchemasDict = engine.SchemasDict{}

	// Get all the files in the models directory
	files, err := os.ReadDir(structsDirPath)

	// If we have an error, we return it
	if err != nil {
		// Send ChanSchemas to channel if we have channel
		if chanSchemas != nil {
			chanSchemas <- engine.ChanSchemas{
				Items: AllSchemasDict,
				Err:   err,
			}
		}

		// Return AllSchemasDict
		return AllSchemasDict, err
	}

	// Loop through all the files
	for _, file := range files {

		// If the file is a directory or is an ignored file we continue
		if file.IsDir() || engine.IsIgnoredFile(file.Name()) {
			continue
		}

		// Extract the schemas from the file
		schemasDict, err := p.extractSchemasDictFromFile(structsDirPath + "/" + file.Name())

		// If we have an error, we return it
		if err != nil {
			// TODO: we should have verbose param to log more
			// log.Println(err)
			continue
		}

		// Append the schemas to the global schemas map
		AllSchemasDict = engine.MergeMaps(AllSchemasDict, schemasDict)
	}
	// Send ChanSchemas to channel if we have channel
	if chanSchemas != nil {
		chanSchemas <- engine.ChanSchemas{
			Items: AllSchemasDict,
			Err:   nil,
		}
	}

	// Return AllSchemasDict
	return AllSchemasDict, nil
}

func (p *openEngine) AddSchemas(schemasDict engine.SchemasDict) OpenEngine {
	p.Components.Schemas = engine.MergeMaps(p.Components.Schemas, schemasDict)
	return p
}

func (p *openEngine) ParseSchemas(baseDirectory string, allIgnoredPaths ...[]string) OpenEngine {
	// Create Mutex for SchemasDict
	mx := &sync.Mutex{}

	// ReInitialize WaitGroup for goroutines
	WgSchema = &sync.WaitGroup{}

	// Create SchemasDict
	AllSchemasDict := p.Components.Schemas

	var ignoredPaths []string
	if len(allIgnoredPaths) > 0 {
		ignoredPaths = allIgnoredPaths[0]
	}
	mergedIgnoredPaths := append(p.GeneralIgnoredPaths, ignoredPaths...)
	// Find all the directories in the baseDirPath
	structsDirectoryPaths, err := engine.FindAllDirectoriesInPath(baseDirectory, &mergedIgnoredPaths)

	// If we have an error, we return it
	if err != nil {
		p.err = err
		return p
	}

	// Create a channel for the schemas
	var chanSchemas = make(chan engine.ChanSchemas, len(structsDirectoryPaths))

	// Close the channel when we are done with it in goroutine
	go func() {
		defer close(chanSchemas)
		WgSchema.Wait()
	}()

	// Loop through all the models paths
	for _, structsDirectoryPath := range structsDirectoryPaths {
		// Add 1 to WaitGroup
		WgSchema.Add(1)

		// Extract the schemas from the directory in a goroutine
		go p.extractSchemasFromDirectory(structsDirectoryPath, chanSchemas)
	}

	// Loop through all the schemas in the channel
	for result := range chanSchemas {
		// Lock the mutex for the SchemasDict
		mx.Lock()

		// Check if we have an error
		if result.Err != nil {
			p.err = result.Err
			return p
		}

		// Merge the schemas to the global schemas map
		AllSchemasDict = engine.MergeMaps(AllSchemasDict, result.Items)

		// Unlock the mutex for the SchemasDict
		mx.Unlock()
	}

	// Return the global schemas map
	p.Components.Schemas = AllSchemasDict
	return p
}
