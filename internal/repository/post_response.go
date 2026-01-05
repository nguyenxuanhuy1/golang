package repository

import "time"

// Struct dùng để trả về JSON gọn gàng
type PostResponse struct {
	ID          int64      `json:"id"`
	ImageID     *int64     `json:"imageId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Topic       string     `json:"topic"`
	Prompt      *string    `json:"prompt,omitempty"`
	HotLevel    int8       `json:"hotLevel"`
	HotAt       *time.Time `json:"hotAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Hàm chuyển Post -> PostResponse
func ToPostResponse(post *Post) PostResponse {
	var prompt *string
	if post.Prompt.Valid {
		prompt = &post.Prompt.String
	}

	var imageID *int64
	if post.ImageID.Valid {
		imageID = &post.ImageID.Int64
	}

	var hotAt *time.Time
	if post.HotAt.Valid {
		hotAt = &post.HotAt.Time
	}

	return PostResponse{
		ID:          post.ID,
		ImageID:     imageID,
		Name:        post.Name,
		Description: post.Description,
		Topic:       post.Topic,
		Prompt:      prompt,
		HotLevel:    post.HotLevel,
		HotAt:       hotAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}
