package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-very-secret-key") // ต้องตรงกับใน Middleware

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login godoc
// @Summary      User Login
// @Description  ตรวจสอบ Username/Password และส่งกลับ JWT Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body      api.LoginRequest  true  "Login Credentials"
// @Success      200    {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. ตรวจสอบ User (ตัวอย่างนี้เป็น Hardcode แต่ในงานจริงต้องเช็คจาก Database)
	// ปกติควรดึง User จาก DB มาเทียบ Password (Bcrypt) และเอา GroupID มาด้วย
	if req.Username == "admin" && req.Password == "password123" {

		// 2. สร้าง Claims (ข้อมูลที่จะใส่ใน Token)
		claims := jwt.MapClaims{
			"user_id":  "U001",
			"group_id": "GRP_ADMIN",                           // สำคัญสำหรับตาราง policy_
			"exp":      time.Now().Add(time.Hour * 24).Unix(), // หมดอายุใน 24 ชม.
		}

		// 3. สร้าง Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}
