package QuestionController

import (
	"net/http"
	"strconv"
	"strings"

	"quiz-project/src/models"

	"github.com/gin-gonic/gin"
)

func CreateQuestions(c *gin.Context) {
	quizID, _ := strconv.ParseInt(c.Param("quizID"), 10, 64)

	// Membuat 10 pertanyaan
	for i := 0; i < 10; i++ {
		var question models.Pertanyaan
		question.Pertanyaan = "Pertanyaan " + strconv.Itoa(i+1)
		question.QuizID = int(quizID)
		question.JawabanBenar = i % 4

		// Opsi Jawaban untuk pertanyaan
		opsiJawaban := []string{"", "", "", ""}
		question.OpsiJawaban = strings.Join(opsiJawaban, ", ")

		models.DB.Create(&question)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Questions created successfully"})
}
