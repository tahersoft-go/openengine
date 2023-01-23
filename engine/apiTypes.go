package engine

// TODO: Check Formats : specially Objects
func OpenAPIFormats(t string) string {
	switch t {
	case "string":
		return "string"
	case "int":
		return "int32"
	case "uint":
		return "int32"
	case "int64":
		return "int64"
	case "uint64":
		return "int64"
	case "float64":
		return "double"
	case "bool":
		return "boolean"
	case "[]string":
		return "string"
	case "[]int":
		return "int32"
	case "[]uint":
		return "int32"
	case "[]int64":
		return "int64"
	case "[]uint64":
		return "int64"
	case "[]float64":
		return "double"
	case "[]bool":
		return "boolean"
	case "[]interface{}":
		return "object"
	case "interface{}":
		return "object"
	case "map[string]interface{}":
		return "object"
	}
	return "string"
}

func OpenAPITypes(t string) string {
	switch t {
	case "string":
		return "string"
	case "int":
		return "integer"
	case "uint":
		return "integer"
	case "int64":
		return "integer"
	case "uint64":
		return "integer"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	case "[]string":
		return "array"
	case "[]int":
		return "array"
	case "[]uint":
		return "array"
	case "[]int64":
		return "array"
	case "[]uint64":
		return "array"
	case "[]float64":
		return "array"
	case "[]bool":
		return "array"
	case "[]interface{}":
		return "array"
	case "interface{}":
		return "object"
	case "map[string]interface{}":
		return "object"
	}
	return "object"
}
