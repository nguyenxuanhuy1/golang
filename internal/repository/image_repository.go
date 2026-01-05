package repository

import (
	"database/sql"
	"time"
)

type Image struct {
	ID          int64
	ImageURL    string
	BlurURL     string
	TinyBlurURL string
	PublicID    string
	ImageType   string
	OwnerID     sql.NullInt64
	CreatedAt   time.Time
}
type ImageRepo interface {
	Create(img *Image) error
	GetByID(id int64) (*Image, error)
	DeleteByID(id int64) error
}

type imageRepo struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepo {
	return &imageRepo{db: db}
}

func (r *imageRepo) Create(img *Image) error {
	query := `
		INSERT INTO images (image_url, blur_url, tiny_blur_url, public_id, image_type, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(
		query,
		img.ImageURL, img.BlurURL, img.TinyBlurURL, img.PublicID, img.ImageType, img.OwnerID, time.Now(),
	).Scan(&id)
	if err != nil {
		return err
	}
	img.ID = id
	return nil
}
func (r *imageRepo) DeleteByID(id int64) error {
	_, err := r.db.Exec("DELETE FROM images WHERE id = $1", id)
	return err
}
func (r *imageRepo) GetByID(id int64) (*Image, error) {
	query := `
		SELECT
			id,
			image_url,
			blur_url,
			tiny_blur_url,
			public_id,
			image_type,
			owner_id,
			created_at
		FROM images
		WHERE id = $1
	`

	var img Image
	err := r.db.QueryRow(query, id).Scan(
		&img.ID,
		&img.ImageURL,
		&img.BlurURL,
		&img.TinyBlurURL,
		&img.PublicID,
		&img.ImageType,
		&img.OwnerID,
		&img.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &img, nil
}
