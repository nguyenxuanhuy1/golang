package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"traingolang/internal/helper"
	"traingolang/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type PostRequest struct {
	Name        string `json:"name"`        // tên bài viết
	Description string `json:"description"` // mô tả
	Topic       string `json:"topic"`       // chủ đề
	Prompt      string `json:"prompt"`      // prompt nếu có
	HotLevel    int8   `json:"hotLevel"`
}
type SearchPostRequest struct {
	Name     string `json:"name"`
	Topic    string `json:"topic"`
	HotLevel *int8  `json:"hotLevel"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
type PostOption struct {
	Value int64  `json:"value"`
	Label string `json:"label"`
}

func CreatePost(postRepo repository.PostRepo, imageRepo repository.ImageRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy JSON data
		dataStr := c.PostForm("data")
		if dataStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data is required"})
			return
		}

		var req PostRequest
		if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data object"})
			return
		}
		// exists, err := postRepo.ExistsByTopic(req.Topic, nil)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }
		// if exists {
		// 	c.JSON(http.StatusConflict, gin.H{"error": "Topic đã tồn tại"})
		// 	return
		// }

		// 2. Upload ảnh + lưu DB dùng helper
		img, err := UploadAndSaveImage(c, "post")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 3. Tạo post
		post := &repository.Post{
			ImageID:     sql.NullInt64{Int64: img.ID, Valid: true},
			Name:        req.Name,
			Description: req.Description,
			Topic:       req.Topic,
			Prompt:      sql.NullString{String: req.Prompt, Valid: req.Prompt != ""},
			HotLevel:    req.HotLevel,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		postID, err := postRepo.CreatePost(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		post.ID = postID

		c.JSON(http.StatusCreated, gin.H{
			"message": "Tạo mới thành công",
		})
	}
}
func SearchPostsHandler(postRepo repository.PostRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SearchPostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		result, err := postRepo.SearchPosts(req.Name, req.Topic, req.HotLevel, req.Page, req.PageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
func UpdatePost(postRepo repository.PostRepo, imageRepo repository.ImageRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		dataStr := c.PostForm("data")
		if dataStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data is required"})
			return
		}

		var req PostRequest
		if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data object"})
			return
		}

		post, err := postRepo.GetByID(postID)
		if err != nil || post == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		// check trùng topic (chỉ khi đổi topic)
		// if req.Topic != post.Topic {
		// 	exists, err := postRepo.ExistsByTopic(req.Topic, &postID)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	if exists {
		// 		c.JSON(http.StatusConflict, gin.H{"error": "Topic đã tồn tại"})
		// 		return
		// 	}
		// }

		// Upload ảnh mới (nếu có)
		if _, err := c.FormFile("image"); err == nil {

			newImg, err := UploadAndSaveImage(c, "post")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if post.ImageID.Valid {
				oldImg, _ := imageRepo.GetByID(post.ImageID.Int64)
				if oldImg != nil {
					_ = helper.DeleteImageFromCloud(oldImg.PublicID)
					_ = imageRepo.DeleteByID(oldImg.ID)
				}
			}

			post.ImageID = sql.NullInt64{Int64: newImg.ID, Valid: true}
		}

		post.Name = req.Name
		post.Description = req.Description
		post.Topic = req.Topic
		post.Prompt = sql.NullString{String: req.Prompt, Valid: req.Prompt != ""}
		post.HotLevel = req.HotLevel
		post.UpdatedAt = time.Now()

		err = postRepo.UpdatePost(post)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Topic đã tồn tại"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công"})
	}
}

func GetPostOptionsHandler(postRepo repository.PostRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := postRepo.GetPostOptions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func DeletePost(postRepo repository.PostRepo, imageRepo repository.ImageRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		post, err := postRepo.GetByID(postID)
		if err != nil || post == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		if post.ImageID.Valid {
			img, _ := imageRepo.GetByID(post.ImageID.Int64)
			if img != nil {
				_ = helper.DeleteImageFromCloud(img.PublicID)
				_ = imageRepo.DeleteByID(img.ID)
			}
		}

		_ = postRepo.DeletePost(postID)

		c.JSON(http.StatusOK, gin.H{"message": "Xoá thành công"})
	}
}
