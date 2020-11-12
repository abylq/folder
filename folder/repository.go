package folder

import (
	"context"
	"github.com/abylq/folder/models"
)

type Repository interface {
	CreateFolder(ctx context.Context, user *models.User, bm *models.Folder) error
}
