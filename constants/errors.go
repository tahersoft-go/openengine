package constants

import "gitlab.hoitek.fi/openapi/openengine/engine"

var DefaultErrorResponses = engine.ErrorResponses{
	"400": engine.Response{
		Description: "Bad Request",
	},
	"401": engine.Response{
		Description: "Unauthorized",
	},
	"403": engine.Response{
		Description: "Forbidden",
	},
	"404": engine.Response{
		Description: "Not Found",
	},
	"405": engine.Response{
		Description: "Method Not Allowed",
	},
	"406": engine.Response{
		Description: "Not Acceptable",
	},
	"408": engine.Response{
		Description: "Request Timeout",
	},
	"409": engine.Response{
		Description: "Conflict",
	},
	"410": engine.Response{
		Description: "Gone",
	},
	"415": engine.Response{
		Description: "Unsupported Media Type",
	},
	"429": engine.Response{
		Description: "Too Many Requests",
	},
	"500": engine.Response{
		Description: "Internal Server Error",
	},
	"501": engine.Response{
		Description: "Not Implemented",
	},
	"502": engine.Response{
		Description: "Bad Gateway",
	},
	"503": engine.Response{
		Description: "Service Unavailable",
	},
	"504": engine.Response{
		Description: "Gateway Timeout",
	},
	"505": engine.Response{
		Description: "HTTP Version Not Supported",
	},

	"422": engine.Response{
		Description: "Unprocessable Entity",
	},
	"423": engine.Response{
		Description: "Locked",
	},
	"424": engine.Response{
		Description: "Failed Dependency",
	},
	"426": engine.Response{
		Description: "Upgrade Required",
	},
	"428": engine.Response{
		Description: "Precondition Required",
	},
	"431": engine.Response{
		Description: "Request Header Fields Too Large",
	},
	"451": engine.Response{
		Description: "Unavailable For Legal Reasons",
	},
	"507": engine.Response{
		Description: "Insufficient Storage",
	},
	"511": engine.Response{
		Description: "Network Authentication Required",
	},

	"520": engine.Response{
		Description: "Unknown Error",
	},
	"521": engine.Response{
		Description: "Web Server Is Down",
	},
	"522": engine.Response{
		Description: "Connection Timed Out",
	},
	"523": engine.Response{
		Description: "Origin Is Unreachable",
	},
	"524": engine.Response{
		Description: "A Timeout Occurred",
	},
	"525": engine.Response{
		Description: "SSL Handshake Failed",
	},
	"526": engine.Response{
		Description: "Invalid SSL Certificate",
	},
	"527": engine.Response{
		Description: "Railgun Error",
	},

	"530": engine.Response{
		Description: "Site is frozen",
	},

	"598": engine.Response{
		Description: "Network read timeout error",
	},
	"599": engine.Response{
		Description: "Network connect timeout error",
	},
}
