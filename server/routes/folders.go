package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
)

func InitFolders(router *gin.RouterGroup) {
	folderRouter := router.Group("/folders")
	folderRouter.Use(middleware.AuthRequired())

	folderRouter.POST("", CreateFolder)
	folderRouter.GET("/root", GetRootFolder)
	folderRouter.GET("/:folderId", GetFolder)
	folderRouter.PATCH("/:folderId", UpdateFolder)
	folderRouter.DELETE("/:folderId", DeleteFolder)
	folderRouter.POST("/:folderId/add-folder", AddFolderToFolder)
	folderRouter.POST("/:folderId/add-event", AddEventToFolder)
}

func GetRootFolder(c *gin.Context) {
	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	folders, err := db.GetChildFolders(nil, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folders"})
		return
	}

	events, err := db.GetEventsInFolder(nil, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	rootFolder := models.Folder{
		Id:       primitive.NilObjectID,
		UserId:   userId,
		ParentId: nil,
		Name:     "Root",
		Folders:  folders,
		Events:   events,
	}

	c.JSON(http.StatusOK, rootFolder)
}

func GetFolder(c *gin.Context) {
	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	folderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	folder, err := db.GetFolderById(folderId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	childFolders, err := db.GetChildFolders(&folderId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get child folders"})
		return
	}
	folder.Folders = childFolders

	events, err := db.GetEventsInFolder(&folderId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events in folder"})
		return
	}
	folder.Events = events

	c.JSON(http.StatusOK, folder)
}

func CreateFolder(c *gin.Context) {
	var body struct {
		Name     string  `json:"name" binding:"required"`
		ParentId *string `json:"parentId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var parentId *primitive.ObjectID
	if body.ParentId != nil {
		pId, err := primitive.ObjectIDFromHex(*body.ParentId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent ID"})
			return
		}
		parentId = &pId
	}

	folder := models.Folder{
		UserId:   userId,
		Name:     body.Name,
		ParentId: parentId,
	}

	id, err := db.CreateFolder(&folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

func UpdateFolder(c *gin.Context) {
	folderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}
	var body struct {
		Name     *string `json:"name"`
		ParentId *string `json:"parentId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updates := bson.M{}
	if body.Name != nil {
		updates["name"] = body.Name
	}
	if body.ParentId != nil {
		if *body.ParentId == "" {
			updates["parentId"] = nil
		} else {
			parentId, err := primitive.ObjectIDFromHex(*body.ParentId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent ID"})
				return
			}
			updates["parentId"] = parentId
		}
	}

	err = db.UpdateFolder(folderId, userId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.Status(http.StatusOK)
}

func AddFolderToFolder(c *gin.Context) {
	parentFolderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent folder ID"})
		return
	}
	var body struct {
		FolderID string `json:"folderId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	childFolderId, err := primitive.ObjectIDFromHex(body.FolderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child folder ID"})
		return
	}

	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updates := bson.M{"parentId": parentFolderId}
	err = db.UpdateFolder(childFolderId, userId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move folder"})
		return
	}

	c.Status(http.StatusOK)
}

func AddEventToFolder(c *gin.Context) {
	folderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}
	var body struct {
		EventID string `json:"eventId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventId, err := primitive.ObjectIDFromHex(body.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = db.SetEventFolder(eventId, &folderId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add event to folder"})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteFolder(c *gin.Context) {
	folderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}
	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = db.DeleteFolder(folderId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}
	c.Status(http.StatusOK)
}
