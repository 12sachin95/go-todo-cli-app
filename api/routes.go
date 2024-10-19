package api

import (
	"go-todo-cli/db"
	"go-todo-cli/models"
	"go-todo-cli/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	userRoutes.POST("/register", register)
	userRoutes.POST("/login", login)
	userRoutes.POST("/logout", logout)
}

func TodoRoutes(router *gin.RouterGroup) {
	protected := router.Group("/todos")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/", getAllTodos)
		protected.GET("/:id", getTodo)
		protected.PUT("/:id", updateTodo)
		protected.POST("/", createTodo)
		protected.DELETE("/:id", deleteTodo)
	}

}

func StartServer() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the environment variables
	port := os.Getenv("PORT")
	uri := os.Getenv("MONGODB_URI")
	if port == "" {
		port = "8080" // Default port if not set
	}

	db.ConnectMongoDB(uri)
	r := gin.Default()

	v1 := r.Group("/todo-app/api/v1")

	// grouping all routes with api/v1
	{
		AuthRoutes(v1)
		TodoRoutes(v1)
	}

	// Start the server on port 8080
	r.Run(":" + port)
}

func register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := services.RegisterUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	token, err := services.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func logout(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	token = strings.Split(token, "Bearer ")[1]

	err := services.LogoutUser(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func getAllTodos(c *gin.Context) {
	todos, _ := services.GetTodos() // Get todos from service layer
	c.JSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {
	idStr := c.Param("id")
	todos, err2 := services.GetTodoByID(idStr) // Get todos from service layer
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	c.JSON(http.StatusOK, todos)
}
func createTodo(c *gin.Context) {
	var newTodo struct {
		Title string `json:"title" binding:"required"` // Require the title field
	}
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	todoToAdd := models.Todo{Title: newTodo.Title}
	result, err := services.AddTodo(todoToAdd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func updateTodo(c *gin.Context) {
	var newTodo models.Todo
	idStr := c.Param("id")
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	result, err := services.UpdateTodo(idStr, newTodo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	_, err := services.DeleteTodo(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
