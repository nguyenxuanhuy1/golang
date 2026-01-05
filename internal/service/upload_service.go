package service

import (
	"context"
	"mime/multipart"

	"traingolang/internal/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadResult struct {
	Image    string
	Blur     string
	TinyBlur string
	PublicID string
}

func UploadImageWithBlur(
	file multipart.File,
	folder string,
) (*UploadResult, error) {

	cld := config.GetCloudinary()
	ctx := context.Background()

	//  Upload ảnh gốc – KHÔNG transformation
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return nil, err
	}

	img, err := cld.Image(resp.PublicID)
	if err != nil {
		return nil, err
	}

	//  Ảnh cực mờ (placeholder)
	img.Transformation = "w_30,q_1,e_blur:1200"
	tinyBlur, err := img.String()
	if err != nil {
		return nil, err
	}

	//  Ảnh mờ (preview)
	img.Transformation = "w_200,q_20,e_blur:400"
	blur, err := img.String()
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Image:    resp.SecureURL,
		Blur:     blur,
		TinyBlur: tinyBlur,
		PublicID: resp.PublicID,
	}, nil
}
