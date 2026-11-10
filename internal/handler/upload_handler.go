package handler

import (
	"api-template/internal/model"
	"api-template/internal/service"
	"api-template/pkg/apperr"
	"api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageSvc *service.StorageService
}

func NewUploadHandler(storageSvc *service.StorageService) *UploadHandler {
	return &UploadHandler{storageSvc: storageSvc}
}

// GeneratePresignedURL godoc
// @Tags upload
// @Summary Generate Presigned Upload URL for MinIO / S3
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.PresignedUploadRequest true "Presigned Upload Request"
// @Success 200 {object} response.Response{data=model.PresignedUploadResponse}
// @Router /api/upload/presigned-url [post]
func (h *UploadHandler) GeneratePresignedURL(c *gin.Context) {
	if h.storageSvc == nil {
		c.Error(apperr.New(apperr.ErrInternal, "storage service is not configured"))
		return
	}

	var req model.PresignedUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.New(apperr.ErrInvalidJSON, err.Error()))
		return
	}

	res, err := h.storageSvc.GeneratePresignedUploadURL(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, res)
}
