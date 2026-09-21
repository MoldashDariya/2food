package docs

import _ "embed"

//go:embed swagger.json
var SwaggerJSON []byte

//go:embed scalar.html
var ScalarHTML []byte // scalar for beautiful API documentation
