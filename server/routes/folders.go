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

	folderRouter.GET("", GetFolders)
	folderRouter.POST("", CreateFolder)
	folderRouter.PATCH("/:folderId", UpdateFolder)
	folderRouter.DELETE("/:folderId", DeleteFolder)
	folderRouter.POST("/:folderId/add-folder", AddFolderToFolder)
	folderRouter.POST("/:folderId/add-event", AddEventToFolder)
}

func GetFolders(c *gin.Context) {
	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	folders, err := db.GetTopLevelFolders(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folders"})
		return
	}

	c.JSON(http.StatusOK, folders)
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

	folder := models.Folder{
		UserId:   userId,
		Name:     body.Name,
		ParentId: body.ParentId,
		Folders:  []primitive.ObjectID{},
		Events:   []primitive.ObjectID{},
	}

	id, err := db.CreateFolder(&folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	if body.ParentId != nil {
		err = db.AddFolderToFolder(*body.ParentId, userId, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add folder to parent"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func UpdateFolder(c *gin.Context) {
	folderId := c.Param("folderId")
	var body struct {
		Name string `json:"name" binding:"required"`
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

	updates := bson.M{"name": body.Name}

	err = db.UpdateFolder(folderId, userId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.Status(http.StatusOK)
}

func AddFolderToFolder(c *gin.Context) {
	parentFolderId := c.Param("folderId")
	var body struct {
		FolderID string `json:"folderId" binding:"required"`
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

	err = db.AddFolderToFolder(parentFolderId, userId, body.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add folder to folder"})
		return
	}

	err = db.UpdateFolderParent(body.FolderID, userId, parentFolderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder parent"})
		return
	}

	c.Status(http.StatusOK)
}

func AddEventToFolder(c *gin.Context) {
	folderId := c.Param("folderId")
	var body struct {
		EventID string `json:"eventId" binding:"required"`
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

	err = db.AddEventToFolder(folderId, userId, body.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add event to folder"})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteFolder(c *gin.Context) {
	folderId := c.Param("folderId")

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
