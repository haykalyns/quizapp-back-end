package AuthController

import (
	"net/http"
	"quiz-project/src/config"
	"quiz-project/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func Login(c *gin.Context) {
	var userInput models.User
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil data user berdasarkan nama
	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Username atau password salah"})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Cek apakah role pengguna sesuai dengan yang diinputkan
	if user.Role != userInput.Role {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Peran pengguna tidak sesuai"})
		return
	}

	// Cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username atau password salah"})
		return
	}

	// Tentukan role pengguna untuk digunakan dalam pembuatan token
	role := user.Role

	// Proses pembuatan JWT
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "quiz-project",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
		Role: role, // Tambahkan peran pengguna sebagai klaim dalam token
	}

	// Menggunakan algoritma yang akan digunakan untuk sign-in
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Pilih kunci berdasarkan peran pengguna
	var key []byte
	if role == "admin" {
		key = config.AdminJWTKey
	} else {
		key = config.UserJWTKey
	}

	// Sign token menggunakan kunci yang sesuai
	token, err := tokenAlgo.SignedString(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set token ke cookie
	c.SetCookie("token", token, 60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Register(c *gin.Context) {
	// Ambil input json
	var userInput models.User
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password hash"})
		return
	}
	userInput.Password = string(hashPassword)

	// Simpan data ke database
	if err := models.DB.Create(&userInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Logout(c *gin.Context) {
	// Hapus token yang ada di cookie
	c.SetCookie("token", "token", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
