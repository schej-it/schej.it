package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/logger"
	"schej.it/server/models"
)

func CreateFolder(folder *models.Folder) (string, error) {
	result, err := FoldersCollection.InsertOne(context.Background(), folder)
	if err != nil {
		logger.StdErr.Panicln(err)
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetTopLevelFolders(userId primitive.ObjectID) ([]models.Folder, error) {
	var folders []models.Folder
	cursor, err := FoldersCollection.Find(context.Background(), bson.M{
		"userId":   userId,
		"parentId": nil,
		"$or": bson.A{
			bson.M{"isDeleted": bson.M{"$exists": false}},
			bson.M{"isDeleted": false},
		},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

func GetFolderById(folderId string, userId primitive.ObjectID) (*models.Folder, error) {
	objectId, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return nil, err
	}

	var folder models.Folder
	err = FoldersCollection.FindOne(context.Background(), bson.M{
		"$and": bson.A{
			bson.M{"_id": objectId},
			bson.M{"userId": userId},
			bson.M{
				"$or": bson.A{
					bson.M{"isDeleted": bson.M{"$exists": false}},
					bson.M{"isDeleted": false},
				},
			},
		},
	}).Decode(&folder)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func UpdateFolder(folderId string, userId primitive.ObjectID, updates bson.M) error {
	objectId, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return err
	}

	_, err = FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": objectId, "userId": userId}, bson.M{"$set": updates})
	return err
}

func AddFolderToFolder(parentFolderId string, userId primitive.ObjectID, childFolderId string) error {
	parentObjectId, err := primitive.ObjectIDFromHex(parentFolderId)
	if err != nil {
		return err
	}
	childObjectId, err := primitive.ObjectIDFromHex(childFolderId)
	if err != nil {
		return err
	}

	_, err = FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": parentObjectId, "userId": userId}, bson.M{"$addToSet": bson.M{"folders": childObjectId}})
	return err
}

func AddEventToFolder(folderId string, userId primitive.ObjectID, eventId string) error {
	objectId, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return err
	}
	eventIdObj, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return err
	}

	_, err = FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": objectId, "userId": userId}, bson.M{"$addToSet": bson.M{"events": eventIdObj}})
	return err
}

func UpdateFolderParent(folderId string, userId primitive.ObjectID, parentId string) error {
	objectId, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return err
	}
	parentObjectId, err := primitive.ObjectIDFromHex(parentId)
	if err != nil {
		return err
	}

	_, err = FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": objectId, "userId": userId}, bson.M{"$set": bson.M{"parentId": parentObjectId.Hex()}})
	return err
}

func DeleteFolder(folderId string, userId primitive.ObjectID) error {
	objectId, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return err
	}

	_, err = FoldersCollection.UpdateOne(context.Background(), bson.M{"_id": objectId, "userId": userId}, bson.M{"$set": bson.M{"isDeleted": true}})
	return err
}
