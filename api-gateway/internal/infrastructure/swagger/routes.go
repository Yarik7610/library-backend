package swagger

import (
	"github.com/Yarik7610/library-backend/api-gateway/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, swaggerHandler Handler) {
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger-json/doc.json", swaggerHandler.GetMergedDocs)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-json/doc.json")))
}
