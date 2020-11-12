package delivery

import (
	"github.com/abylq/folder/folder"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc folder.UseCase)  {
	h := NewHandler()

	folders := router.Group("/folder")
	{
		folders.POST("",h.Create)
	}
}
