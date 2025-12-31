package handler

import (
	"net/http"

	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username và Password phải tối thiểu 6 ký tự",
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "hash password failed",
		})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)

	_, err = userRepo.Create(req.Username, string(passwordHash))
	if err != nil {
		switch err {
		case repository.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Tài khoản đã tồn tại",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký tài khoản thành công",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)

	user, err := userRepo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	// 1. Check user bị khoá
	if user.Locked {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tài khoản đã bị khoá",
		})
		return
	}

	// 2. Check password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	// 3. Generate access token (15 phút)
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generate access token failed"})
		return
	}

	// 4. Generate refresh token (1 giờ)
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generate refresh token failed"})
		return
	}

	// 5. Response
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Profile(c *gin.Context) {
	// 1. Lấy claims từ middleware
	claimsAny, exists := c.Get(auth.ContextUserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	userID := claims.UserID

	// 2. Query DB
	userRepo := repository.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"avatar":   user.Avatar,
	})
}
