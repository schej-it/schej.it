package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/logger"
	"schej.it/server/models"
)

func CreateFolder(folder *models.Folder) (primitive.ObjectID, error) {
	result, err := FoldersCollection.InsertOne(context.Background(), folder)
	if err != nil {
		logger.StdErr.Panicln(err)
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetFolderById(folderId primitive.ObjectID, userId primitive.ObjectID) (*models.Folder, error) {
	var folder models.Folder
	err := FoldersCollection.FindOne(context.Background(), bson.M{
		"_id":    folderId,
		"userId": userId,
		"$or": bson.A{
			bson.M{"isDeleted": bson.M{"$exists": false}},
			bson.M{"isDeleted": false},
		},
	}).Decode(&folder)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func GetChildFolders(parentFolderId *primitive.ObjectID, userId primitive.ObjectID) ([]models.Folder, error) {
	cursor, err := FoldersCollection.Find(context.Background(), bson.M{
		"parentId": parentFolderId,
		"userId":   userId,
		"$or": bson.A{
			bson.M{"isDeleted": bson.M{"$exists": false}},
			bson.M{"isDeleted": false},
		},
	})
	if err != nil {
		return nil, err
	}

	var childFolders []models.Folder
	if err = cursor.All(context.Background(), &childFolders); err != nil {
		return nil, err
	}

	return childFolders, nil
}

func GetEventsInFolder(folderId *primitive.ObjectID, userId primitive.ObjectID) ([]models.Event, error) {
	cursor, err := EventsCollection.Find(context.Background(), bson.M{
		"folderId": folderId,
		"ownerId":  userId,
		"$or": bson.A{
			bson.M{"isDeleted": bson.M{"$exists": false}},
			bson.M{"isDeleted": false},
		},
	})
	if err != nil {
		return nil, err
	}

	var events []models.Event
	if err = cursor.All(context.Background(), &events); err != nil {
		return nil, err
	}

	return events, nil
}

func UpdateFolder(folderId primitive.ObjectID, userId primitive.ObjectID, updates bson.M) error {
	_, err := FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": folderId, "userId": userId}, bson.M{"$set": updates})
	return err
}

func SetEventFolder(eventId primitive.ObjectID, folderId *primitive.ObjectID, userId primitive.ObjectID) error {
	_, err := EventsCollection.UpdateOne(context.Background(), bson.M{"_id": eventId, "ownerId": userId}, bson.M{"$set": bson.M{"folderId": folderId}})
	return err
}

func DeleteFolder(folderId primitive.ObjectID, userId primitive.ObjectID) error {
	ctx := context.Background()
	// Mark this folder as deleted
	_, err := FoldersCollection.UpdateOne(ctx, bson.M{"_id": folderId, "userId": userId}, bson.M{"$set": bson.M{"isDeleted": true}})
	if err != nil {
		return err
	}

	// Mark all events in this folder as deleted
	_, err = EventsCollection.UpdateMany(ctx, bson.M{"folderId": folderId, "ownerId": userId}, bson.M{"$set": bson.M{"isDeleted": true}})
	if err != nil {
		return err
	}

	// Find all child folders
	cursor, err := FoldersCollection.Find(ctx, bson.M{"parentId": folderId, "userId": userId, "$or": bson.A{
		bson.M{"isDeleted": bson.M{"$exists": false}},
		bson.M{"isDeleted": false},
	}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var childFolders []struct {
		Id primitive.ObjectID `bson:"_id"`
	}
	if err = cursor.All(ctx, &childFolders); err != nil {
		return err
	}

	// Recursively delete all child folders
	for _, child := range childFolders {
		err := DeleteFolder(child.Id, userId)
		if err != nil {
			return err
		}
	}

	return nil
}
