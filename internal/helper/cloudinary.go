package helper

import (
	"context"
	"errors"
	"fmt"

	"traingolang/internal/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func DeleteImageFromCloud(publicID string) error {
	if publicID == "" {
		return errors.New("public_id is empty")
	}

	cld := config.GetCloudinary()

	invalidate := true
	res, err := cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID:     publicID,
			ResourceType: "image",
			Invalidate:   &invalidate,
		},
	)
	if err != nil {
		return fmt.Errorf("cloudinary delete error: %w", err)
	}

	switch res.Result {
	case "ok", "not found":
		return nil
	default:
		return fmt.Errorf("cloudinary delete failed: %s", res.Result)
	}
}
