package handlers

import (
	"crypto/subtle"
	"ninimenu/internal/config"
	"ninimenu/internal/middleware"
	"ninimenu/internal/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

func AppLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "密码不能为空")
		return
	}

	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(config.C.AppPassword)) != 1 {
		utils.Unauthorized(c, "应用密码错误")
		return
	}

	utils.Success(c, gin.H{"verified": true})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "密码不能为空")
		return
	}

	if req.Password != config.C.AdminPassword {
		utils.Unauthorized(c, "密码错误")
		return
	}

	token, err := middleware.GenerateToken()
	if err != nil {
		utils.InternalError(c, "生成Token失败")
		return
	}

	utils.Success(c, gin.H{"token": token})
}
