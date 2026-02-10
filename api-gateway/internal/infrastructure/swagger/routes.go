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
		userDocs := fetchDocsJSON(microservice.USER_ADDRESS)
		catalogDocs := fetchDocsJSON(microservice.CATALOG_ADDRESS)
		subDocs := fetchDocsJSON(microservice.SUBSCRIPTIONS_ADDRESS)

		mergedDoc := mergeDocs(userDocs, catalogDocs, subDocs)
		c.JSON(http.StatusOK, mergedDoc)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-json/doc.json")))
}
