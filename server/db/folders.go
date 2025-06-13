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

func GetAllFolders(userId primitive.ObjectID) ([]models.Folder, error) {
	cursor, err := FoldersCollection.Find(context.Background(), bson.M{
		"userId": userId,
		"$or": bson.A{
			bson.M{"isDeleted": bson.M{"$exists": false}},
			bson.M{"isDeleted": false},
		},
	})
	if err != nil {
		return nil, err
	}

	var folders []models.Folder
	if err = cursor.All(context.Background(), &folders); err != nil {
		return nil, err
	}

	return folders, nil
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
	ctx := context.Background()

	// Remove any existing mapping for this event
	_, err := FolderEventsCollection.DeleteMany(ctx, bson.M{"eventId": eventId, "userId": userId})
	if err != nil {
		return err
	}

	// If folderId is nil, we just un-assign it. Otherwise, create a new mapping
	if folderId != nil {
		_, err = FolderEventsCollection.InsertOne(ctx, models.FolderEvent{
			FolderId: *folderId,
			EventId:  eventId,
			UserId:   userId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteFolder(folderId primitive.ObjectID, userId primitive.ObjectID) error {
	ctx := context.Background()
	// Mark this folder as deleted
	_, err := FoldersCollection.UpdateOne(ctx, bson.M{"_id": folderId, "userId": userId}, bson.M{"$set": bson.M{"isDeleted": true}})
	if err != nil {
		return err
	}

	// Find all event mappings for this folder
	cursor, err := FolderEventsCollection.Find(ctx, bson.M{"folderId": folderId, "userId": userId})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var mappings []models.FolderEvent
	if err = cursor.All(ctx, &mappings); err != nil {
		return err
	}

	eventIds := make([]primitive.ObjectID, len(mappings))
	for i, m := range mappings {
		eventIds[i] = m.EventId
	}

	// Mark all events in this folder as deleted
	if len(eventIds) > 0 {
		_, err = EventsCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": eventIds}, "ownerId": userId}, bson.M{"$set": bson.M{"isDeleted": true}})
		if err != nil {
			return err
		}
	}

	// Delete the mappings
	_, err = FolderEventsCollection.DeleteMany(ctx, bson.M{"folderId": folderId, "userId": userId})
	if err != nil {
		return err
	}

	return nil
}
