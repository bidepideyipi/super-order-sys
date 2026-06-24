package handler

import (
	"net/http"

	"super-order-web/internal/config"
	"super-order-web/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	cfg *config.UserConfig
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(cfg *config.UserConfig) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

// Login 登录
// @Summary 登录
// @Description 用户登录获取 token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, LoginResponse{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	if req.Username != h.cfg.Username || req.Password != h.cfg.Password {
		c.JSON(http.StatusOK, LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	token, err := jwt.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Message: "生成 token 失败",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
	})
}

// Logout 登出
// @Summary 登出
// @Description 用户登出
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// Check 检查认证状态
// @Summary 检查认证
// @Description 检查当前 token 是否有效
// @Tags 认证
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /auth/check [get]
func (h *AuthHandler) Check(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"username": username,
	})
}
