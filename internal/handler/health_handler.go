package handler

import (
	"time"

	"api-template/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

// Check godoc
// @Tags system
// @Summary System Health Check
// @Produce json
// @Success 200 {object} response.Response{data=HealthResponse}
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "disabled"
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			dbStatus = "unhealthy: " + err.Error()
		} else {
			dbStatus = "healthy"
		}
	}

	res := HealthResponse{
		Status:    "running",
		Database:  dbStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	response.Success(c, res)
}

