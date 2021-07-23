# APIKey to JWT Mediation Go Plugin

Built to be run natively as a package by Tyk Gateways.

This plugin is intended to run as the "post" part of the request lifecycle.

The plugin will create a JWT with the subject set to the subject on the developer profile metadata, and let the request continue. Otherwise it will return an error.

# Build the plugin binary (shared object) without using go
`docker run --rm  -v `pwd`:/plugin-source tykio/tyk-plugin-compiler:v3.1.2 api-key-to-jwt.so`


# Build the plugin binary (shared object) using go
In the root of the "api-key-to-jwt.go" file, run

`go build -o ./middleware/go/api-key-to-jwt.so -buildmode=plugin ./middleware/go`


Put the generated so file somewhere Tyk Gateway can access it.

# Setup your API
in API Designer, click on "Raw API Definition"
1. Set ` "driver": "goplugin"`
2. Set the request lifecycle to be run
```
"custom_middleware": {
	"post": [
		{
			"name": "ApiKeyToJwt",
			"path": "/opt/tyk-gateway/api-key-to-jwt.so"
		}
	],
```
"name" has to be the name of the GO function
"path" is wherever you put the binary generated previously
