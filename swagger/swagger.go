package swagger

import "github.com/nanotaboada/go-samples-gin-restful/docs"

func Setup() {
	docs.SwaggerInfo.Title = "go-samples-gin-restful"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Description = "🧪 Proof of Concept for a RESTful API made with Go and Gin"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
