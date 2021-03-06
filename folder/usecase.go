package folder

import (
	"context"
	"github.com/abylq/folder/models"
)

type UseCase interface {
	CreateFolder(ctx context.Context, user *models.User, title string) error
	GetFolders(ctx context.Context, user *models.User) ([]*models.Folder, error)
}
