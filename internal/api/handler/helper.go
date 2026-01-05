package handler // hoặc package service nếu muốn

import (
	"traingolang/internal/config"
	"traingolang/internal/repository"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

func UploadAndSaveImage(c *gin.Context, imageType string) (*repository.Image, error) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, err := service.UploadImageWithBlur(file, imageType)
	if err != nil {
		return nil, err
	}

	imageRepo := repository.NewImageRepository(config.DB)
	img := repository.Image{
		ImageURL:    result.Image,
		BlurURL:     result.Blur,
		TinyBlurURL: result.TinyBlur,
		PublicID:    result.PublicID,
		ImageType:   imageType,
	}

	if err := imageRepo.Create(&img); err != nil {
		return nil, err
	}

	return &img, nil
}
