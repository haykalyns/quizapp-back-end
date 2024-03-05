package QuizController

import (
	"net/http"
	"quiz-project/src/models"
	"time"

	"github.com/gin-gonic/gin"
)

func ListQuiz(c *gin.Context) {
	var quizzes []models.Quiz

	models.DB.Find(&quizzes)

	c.JSON(http.StatusOK, gin.H{"quizzes": quizzes})
}

func CreateQuiz(c *gin.Context) {
	var quiz models.Quiz

	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Parsing waktu mulai
	layout := "2006-01-02 15:04:05" // Format yang diinginkan: YYYY-MM-DD HH:MM:SS
	waktuMulaiString := quiz.WaktuMulai.Format(layout)
	waktuMulai, err := time.Parse(layout, waktuMulaiString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Format waktu mulai tidak valid"})
		return
	}

	// Parsing waktu selesai
	waktuSelesaiString := quiz.WaktuSelesai.Format(layout)
	waktuSelesai, err := time.Parse(layout, waktuSelesaiString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Format waktu selesai tidak valid"})
		return
	}

	// Mengatur lokasi zona waktu menjadi Waktu Indonesia Barat (WIB)
	lokasiWIB, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat lokasi zona waktu"})
		return
	}

	// Konversi waktu ke zona waktu Indonesia Barat (WIB)
	waktuMulai = waktuMulai.In(lokasiWIB)
	waktuSelesai = waktuSelesai.In(lokasiWIB)

	// Menyimpan data ke dalam database
	quiz.WaktuMulai = waktuMulai
	quiz.WaktuSelesai = waktuSelesai

	// Menyimpan data ke dalam database
	if err := models.DB.Create(&quiz).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan data ke dalam database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quiz": quiz})
}

func UpdateQuiz(c *gin.Context) {
	var quiz models.Quiz
	id := c.Param("id")

	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&quiz).Where("id = ?", id).Updates(&quiz).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "quiz tidak dapat diupdate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diperbarui"})
}

func DeleteQuiz(c *gin.Context) {
	var quiz models.Quiz
	id := c.Param("id")

	if models.DB.Delete(&quiz, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat hapus quiz"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus"})
}

func ListAnswers(c *gin.Context) {
	// Implementasi kode untuk menampilkan jawaban peserta berdasarkan ID kuis
}

func StartQuiz(c *gin.Context) {
	// Mendapatkan ID quiz dari parameter URL
	quizID := c.Param("id")

	// Lakukan pengecekan apakah quiz tersedia
	var quiz models.Quiz
	if err := models.DB.First(&quiz.Pertanyaan, quizID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Quiz not found"})
		return
	}

	// Lakukan pengecekan apakah waktu quiz sudah dimulai
	if time.Now().Before(quiz.WaktuMulai) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Quiz has not started yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz started successfully"})
}

// Handler untuk mengirimkan jawaban quiz
func SubmitQuiz(c *gin.Context) {
	// Mendapatkan ID quiz dari parameter URL
	quizID := c.Param("id")

	// Lakukan pengecekan apakah quiz tersedia
	var quiz models.Quiz
	if err := models.DB.First(&quiz.Pertanyaan, quizID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Quiz not found"})
		return
	}

	// Lakukan pengecekan apakah waktu quiz sudah dimulai
	if time.Now().Before(quiz.WaktuSelesai) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Quiz has not started yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz submited successfully"})
}
