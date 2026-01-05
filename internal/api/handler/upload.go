package handler

import (
	"database/sql"
	"net/http"

	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

// phục vụ upload riêng lẻ ảnh
func UploadHandler(c *gin.Context) {
	//  Validate image_type
	imageType := c.PostForm("image_type")
	if imageType == "" {
		c.JSON(400, gin.H{"error": "image_type required"})
		return
	}

	// 2Lấy file
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}
	defer file.Close()

	// 3Upload Cloudinary
	result, err := service.UploadImageWithBlur(file, imageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Lấy user_id từ token
	claimsAny, _ := c.Get(auth.ContextUserKey)
	claims := claimsAny.(*auth.Claims)

	imageRepo := repository.NewImageRepository(config.DB)

	img := repository.Image{
		ImageURL:    result.Image,
		BlurURL:     result.Blur,
		TinyBlurURL: result.TinyBlur,
		PublicID:    result.PublicID,
		ImageType:   imageType,
		OwnerID:     sql.NullInt64{Int64: claims.UserID, Valid: true},
	}

	if err := imageRepo.Create(&img); err != nil {
		c.JSON(500, gin.H{"error": "save image failed"})
		return
	}

	// Trả về thông tin ảnh
	c.JSON(http.StatusOK, gin.H{
		"image_id": img.ID,
		"url":      img.ImageURL,
		"blur_url": img.BlurURL,
		"tiny_url": img.TinyBlurURL,
	})
}
