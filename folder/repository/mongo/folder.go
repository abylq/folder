package mongo

import (
	"context"
	"github.com/abylq/folder/models"
	"go.mongodb.org/mongo-driver/bson"
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
	return &FolderRepository {
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

func (f FolderRepository) GetFolders(ctx context.Context, user *models.User) ([]*models.Folder, error) {
	uid, _ := primitive.ObjectIDFromHex(user.ID)
	cur, err := f.db.Find(ctx, bson.M{
		"userId": uid,
	})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*Folder, 0)

	for cur.Next(ctx) {
		user := new(Folder)
		err := cur.Decode(user)
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return toFolders(out), nil
}

func toModel(f *models.Folder) *Folder {
	uid, _ := primitive.ObjectIDFromHex(f.UserID)

	return &Folder{
		UserID: uid,
		Title: f.Title,
	}
}
func toFolder(b *Folder) *models.Folder {
	return &models.Folder{
		ID:     b.ID.Hex(),
		UserID: b.UserID.Hex(),
		Title:  b.Title,
	}
}

func toFolders(bs []*Folder) []*models.Folder {
	out := make([]*models.Folder, len(bs))

	for i, f := range bs {
		out[i] = toFolder(f)
	}

	return out
}