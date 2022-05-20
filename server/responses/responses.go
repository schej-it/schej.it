package responses

type Error struct {
	Error interface{} `json:"error" binding:"required"`
}
