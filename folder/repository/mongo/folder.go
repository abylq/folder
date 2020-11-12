package mongo

import (
	"context"
	"github.com/abylq/folder/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Folder struct {
	ID 		 primitive.ObjectID `bson:"_id"`
	UserID 	 primitive.ObjectID `bson:"user_id"`
	Title 	 string `bson:"title"`
}

type FolderRepository struct {
	db *mongo.Collection
}

func NewFolderRepository(db *mongo.Database, collection string) *FolderRepository {
	return &FolderRepository{
		db: db.Collection(collection),
	}
}

func (f FolderRepository) CreateFolder(ctx context.Context, user *models.User, fm *models.Folder) error {
	fm.UserID = user.ID
	model := toModel(fm)

	res, err := f.db.InsertOne(ctx,model)
	if err != nil {
		return err
	}

	fm.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func toModel(f *models.Folder) *Folder {
	uid, _ := primitive.ObjectIDFromHex(f.UserID)

	return &Folder{
		UserID: uid,
		Title: f.Title,
	}
}