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

func AuthRoutes(router *gin.Engine) {
	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/logout", logout)
}

func TodoRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/todos", getAllTodos)
		protected.GET("/todos/:id", getTodo)
		protected.PUT("/todos/:id", updateTodo)
		protected.POST("/todos", createTodo)
		protected.DELETE("/todos/:id", deleteTodo)
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

	// // Define routes
	// r.POST("/login", login)       // New login route
	// r.POST("/register", register) // New login route
	// r.POST("/logout", logout)     // New login route

	// protected := r.Group("/")
	// protected.Use(AuthMiddleware())
	// {
	// 	protected.GET("/todos", getAllTodos)
	// 	protected.GET("/todos/:id", getTodo)
	// 	protected.PUT("/todos/:id", updateTodo)
	// 	protected.POST("/todos", createTodo)
	// 	protected.DELETE("/todos/:id", deleteTodo)
	// }

	// group the routes
	AuthRoutes(r)
	TodoRoutes(r)

	// Protect todo routes with middleware

	// Start the server on port 8080
	r.Run(":" + port)
}

// // Register handler
// func register(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	if err := Register(user.Username, user.Password); err != nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "models.User registered successfully"})
// }

// // Login handler
// func login(c *gin.Context) {
// 	var credentials models.User
// 	if err := c.ShouldBindJSON(&credentials); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	token, err := Authenticate(credentials.Username, credentials.Password)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }
// func logout(c *gin.Context) {
// 	// Here, you could add logic to blacklist the token if needed.
// 	// For simplicity, we are just responding with a logout message.
// 	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
// }

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
	// Convert the ID from string to int
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	// 	return
	// }
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
	// Convert the ID from string to int
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	// 	return
	// }
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
	// Convert the ID from string to int
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	// 	return
	// }
	_, err := services.DeleteTodo(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
