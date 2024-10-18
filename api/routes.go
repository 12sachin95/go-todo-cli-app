package api

import (
	"go-todo-cli/todo"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	// Define routes
	r.POST("/login", login)       // New login route
	r.POST("/register", register) // New login route
	r.POST("/logout", logout)     // New login route

	// Protect todo routes with middleware
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/todos", getAllTodos)
		protected.GET("/todos/:id", getTodo)
		protected.PUT("/todos/:id", updateTodo)
		protected.POST("/todos", createTodo)
		protected.DELETE("/todos/:id", deleteTodo)
	}

	// Start the server on port 8080
	r.Run(":8080")
}

// Register handler
func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := Register(user.Username, user.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handler
func login(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func logout(c *gin.Context) {
	// Here, you could add logic to blacklist the token if needed.
	// For simplicity, we are just responding with a logout message.
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
}
func getAllTodos(c *gin.Context) {
	todos := todo.GetTodos() // Get todos from service layer
	c.JSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {
	idStr := c.Param("id")
	// Convert the ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	todos, err2 := todo.GetTodoByID(id) // Get todos from service layer
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
	todoToAdd := todo.Todo{Title: newTodo.Title}
	todo.AddTodo(todoToAdd)
	c.JSON(http.StatusOK, newTodo)
}

func updateTodo(c *gin.Context) {
	var newTodo todo.Todo
	idStr := c.Param("id")
	// Convert the ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	todo.UpdateTodoByID(id, newTodo)
	c.JSON(http.StatusOK, newTodo)
}

func deleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	// Convert the ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = todo.DeleteTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
