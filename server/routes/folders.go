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
	folderRouter := router.Group("/user/folders")
	folderRouter.Use(middleware.AuthRequired())

	folderRouter.GET("", GetAllFolders)
	folderRouter.POST("", CreateFolder)
	folderRouter.GET("/:folderId", GetFolder)
	folderRouter.PATCH("/:folderId", UpdateFolder)
	folderRouter.DELETE("/:folderId", DeleteFolder)
}

// @Summary Get all folders
// @Tags folders
// @Produce json
// @Success 200 {array} models.Folder "A list of all folders for the user"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Failed to get folders"
// @Router /user/folders [get]
func GetAllFolders(c *gin.Context) {
	session := sessions.Default(c)
	userIdString := session.Get("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	folders, err := db.GetAllFolders(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folders"})
		return
	}

	if folders == nil {
		c.JSON(http.StatusOK, []models.Folder{})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// @Summary Get a folder by its ID and its contents
// @Tags folders
// @Produce json
// @Param folderId path string true "Folder ID"
// @Success 200 {object} models.Folder "The folder object with events"
// @Failure 400 {object} map[string]string "Invalid user ID or folder ID"
// @Failure 404 {object} map[string]string "Folder not found"
// @Failure 500 {object} map[string]string "Failed to get events in folder"
// @Router /user/folders/{folderId} [get]
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

	events, err := db.GetEventsInFolder(folderId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events in folder"})
		return
	}
	folder.EventIds = events

	c.JSON(http.StatusOK, folder)
}

type CreateFolderResponse struct {
	Id string `json:"id"`
}

// @Summary Create a new folder
// @Tags folders
// @Accept json
// @Produce json
// @Param payload body object{name=string,color=string} true "Folder name and optional color"
// @Success 201 {object} CreateFolderResponse "The ID of the created folder"
// @Failure 400 {object} map[string]string "Invalid user ID or request body"
// @Failure 500 {object} map[string]string "Failed to create folder"
// @Router /user/folders [post]
func CreateFolder(c *gin.Context) {
	var body struct {
		Name  string  `json:"name" binding:"required"`
		Color *string `json:"color"`
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
		UserId: userId,
		Name:   body.Name,
		Color:  body.Color,
	}

	id, err := db.CreateFolder(&folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, CreateFolderResponse{Id: id.Hex()})
}

// @Summary Update a folder's name or color
// @Tags folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param payload body object{name=string,color=string} true "New folder name and/or color"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid user ID or folder ID"
// @Failure 500 {object} map[string]string "Failed to update folder"
// @Router /user/folders/{folderId} [patch]
func UpdateFolder(c *gin.Context) {
	folderId, err := primitive.ObjectIDFromHex(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}
	var body struct {
		Name  *string `json:"name"`
		Color *string `json:"color"`
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
	if body.Color != nil {
		updates["color"] = body.Color
	}

	err = db.UpdateFolder(folderId, userId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete a folder
// @Tags folders
// @Produce json
// @Param folderId path string true "Folder ID"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid user ID or folder ID"
// @Failure 500 {object} map[string]string "Failed to delete folder"
// @Router /user/folders/{folderId} [delete]
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
