package http

import (
	"github.com/abylq/folder/auth"
	"github.com/abylq/folder/folder"
	"github.com/abylq/folder/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Folder struct {
	ID 	  string `json:"id"`
	Title string `json:"title"`
}

type Handler struct {
	useCase folder.UseCase
}

func NewHandler(useCase folder.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Title string `json:"title"`
}

type getResponse struct {
	Folders []*Folder `json:"folders"`
}

func (h *Handler) Create(c *gin.Context) {
	input := new(createInput)
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.CreateFolder(c.Request.Context(), user, input.Title); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Get(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	bms, err := h.useCase.GetFolders(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		Folders: toFolders(bms),
	})
}