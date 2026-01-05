package config

import (
	"log"
	"os"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
)

var (
	cld  *cloudinary.Cloudinary
	once sync.Once
)

func GetCloudinary() *cloudinary.Cloudinary {
	once.Do(func() {
		c, err := cloudinary.NewFromParams(
			os.Getenv("CLOUD_NAME"),
			os.Getenv("CLOUD_API_KEY"),
			os.Getenv("CLOUD_API_SECRET"),
		)
		if err != nil {
			log.Fatal("Cloudinary init error:", err)
		}
		cld = c
	})
	return cld
}
