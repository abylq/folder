package folder

import (
	"context"
	"github.com/abylq/folder/models"
)

type Repository interface {
	CreateFolder(ctx context.Context, user *models.User, bm *models.Folder) error
	GetFolders(ctx context.Context, user *models.User) ([]*models.Folder, error)
	DeleteFolder(ctx context.Context, user *models.User, id string) error
}
