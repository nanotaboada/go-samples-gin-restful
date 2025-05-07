// Package swagger registers Swagger UI and routes for API documentation.
package swagger

import "github.com/nanotaboada/go-samples-gin-restful/docs"

// Setup configures the Swagger UI and routes for API documentation.
func Setup() {
	docs.SwaggerInfo.Title = "go-samples-gin-restful"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Description = "🧪 Proof of Concept for a RESTful API made with Go and Gin"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
