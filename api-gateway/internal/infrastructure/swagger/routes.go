package swagger

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend/api-gateway/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger-json/doc.json", func(c *gin.Context) {
		userDoc := fetchDocsJSON(microservice.USER_ADDRESS)
		catalogDoc := fetchDocsJSON(microservice.USER_ADDRESS)
		subDoc := fetchDocsJSON(microservice.USER_ADDRESS)

		mergedDoc := mergeDocs(userDoc, catalogDoc, subDoc)
		c.JSON(http.StatusOK, mergedDoc)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-json/doc.json")))
}
