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
	HotLevel *int8  `json:"hotLevel"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
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

		result, err := postRepo.SearchPosts(req.Name, req.HotLevel, req.Page, req.PageSize)
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
		var req PostRequest
		_ = json.Unmarshal([]byte(dataStr), &req)

		post, err := postRepo.GetByID(postID)
		if err != nil || post == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		if _, err := c.FormFile("image"); err == nil {

			// 1. Upload ảnh mới trước
			newImg, err := UploadAndSaveImage(c, "post")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// 2. Xoá ảnh cũ (nếu có)
			if post.ImageID.Valid {
				oldImg, _ := imageRepo.GetByID(post.ImageID.Int64)
				if oldImg != nil {
					_ = helper.DeleteImageFromCloud(oldImg.PublicID)
					_ = imageRepo.DeleteByID(oldImg.ID)
				}
			}

			// 3. Gán image mới
			post.ImageID = sql.NullInt64{Int64: newImg.ID, Valid: true}
		}

		post.Name = req.Name
		post.Description = req.Description
		post.Topic = req.Topic
		post.Prompt = sql.NullString{String: req.Prompt, Valid: req.Prompt != ""}
		post.HotLevel = req.HotLevel
		post.UpdatedAt = time.Now()

		_ = postRepo.UpdatePost(post)

		c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công"})
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
