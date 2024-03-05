package main

import (
	// "log"
	// "net/http"

	"quiz-project/src/controllers/AuthController"
	"quiz-project/src/controllers/QuestionController"
	"quiz-project/src/controllers/QuizController"

	"quiz-project/src/middleware"
	"quiz-project/src/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	r.POST("/login", AuthController.Login)
	r.POST("/register", AuthController.Register)
	r.GET("/logout", AuthController.Logout)
	r.GET("/user", AuthController.GetAllUsers)
	// r.GET("/quiz", QuizController.ListQuiz)
	r.POST("/quiz", QuizController.CreateQuiz)
	r.POST("/create-questions/:quizID", QuestionController.CreateQuestions)

	// ADMIN
	admin := r.Group("/admin")
	admin.Use(middleware.AuthenticateAdmin())
	{
		// CRUD Quiz By Admin
		admin.GET("/quiz", QuizController.ListQuiz)
		admin.POST("/quiz", QuizController.CreateQuiz)
		admin.PUT("/quiz/:id", QuizController.UpdateQuiz)
		admin.DELETE("/quiz/:id", QuizController.DeleteQuiz)
		admin.GET("/answers/:quizID", QuizController.ListAnswers)
		admin.POST("/create-questions/:quizID", QuestionController.CreateQuestions)
	}

	// USER
	user := r.Group("/user")
	user.Use(middleware.AuthenticateUser())
	{
		user.GET("/quiz", QuizController.ListQuiz)
		user.POST("/quiz/:id/start", QuizController.StartQuiz)
		user.POST("/quiz/:id/submit", QuizController.SubmitQuiz)
	}

	// // Evaluasi
	// r.GET("/api/scores/:quizID", middleware.AuthenticateAdmin(), QuizController.CalculateScores)

	// Run the server on port 8082
	r.Run(":8082")
}
