package delivery

import (
	"github.com/abylq/folder/folder"
	"github.com/abylq/folder/folder/delivery/http"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc folder.UseCase)  {
	h := http.NewHandler(uc)


	folders := router.Group("/folder")
	{
		folders.POST("",h.Create)
	}

}
