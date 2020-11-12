package usecase

import (
	"context"
	"github.com/abylq/folder/folder"
	"github.com/abylq/folder/models"
)

type FolderUseCase struct {
	folderRepo folder.Repository
}

func NewFolderUseCase(folderRepo folder.Repository) *FolderUseCase {
	return &FolderUseCase{
		folderRepo: folderRepo,
	}
}

func (f FolderUseCase) CreateFolder(ctx context.Context, user *models.User, url, title string) error {
	fm := &models.Folder{
		Title: title,
	}

	return f.folderRepo.CreateFolder(ctx,user,fm)
}