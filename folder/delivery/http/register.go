package http

import (
	"github.com/abylq/folder/folder"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc folder.UseCase)  {
	h := NewHandler(uc)


	folders := router.Group("/folder")
	{
		folders.POST("",h.Create)
		folders.GET("", h.Get)
	}

}
