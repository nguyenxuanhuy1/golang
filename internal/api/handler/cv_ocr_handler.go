package handler

import (
	"io"
	"net/http"
	"strings"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxFileSize = 5 << 20

type CVOCRHandler struct {
	cvService *service.CVOCRService
}

func NewCVOCRHandler() *CVOCRHandler {
	return &CVOCRHandler{
		cvService: service.NewCVOCRService(),
	}
}

func (h *CVOCRHandler) UploadCV(c *gin.Context) {
	requestId := uuid.New().String()

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		return
	}

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read file failed"})
		return
	}

	text, err := h.cvService.ReadCVImage(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requestId": requestId,
		"text":      text,
	})
}
