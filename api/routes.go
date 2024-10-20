package api

import (
	"go-todo-cli/db"
	"go-todo-cli/models"
	"go-todo-cli/services"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Initialize a validator instance
var validate *validator.Validate

func AuthRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	userRoutes.POST("/register", register)
	userRoutes.POST("/login", login)
	userRoutes.POST("/logout", logout)
	userRoutes.GET("/details/:id", getUserDetails)
}

func TodoRoutes(router *gin.RouterGroup) {
	protected := router.Group("/todos")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/", ExtractUserIDFromJWT, getAllTodos)
		protected.GET("/:id", ExtractUserIDFromJWT, getTodo)
		protected.PUT("/:id", ExtractUserIDFromJWT, updateTodo)
		protected.POST("/", ExtractUserIDFromJWT, createTodo)
		protected.DELETE("/:id", ExtractUserIDFromJWT, deleteTodo)
	}

}

func StartServer() {
	validate = validator.New()

	// Get the environment variables
	port := os.Getenv("PORT")
	uri := os.Getenv("MONGODB_URI")
	if port == "" {
		port = "8000" // Default port if not set
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
		Username string `bson:"username" json:"username" validate:"required,min=3,max=32"`
		Email    string `bson:"email" json:"email" validate:"required,email"`
		Password string `bson:"password" json:"password" validate:"required,min=3"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Validate the user struct
	if err := validate.Struct(&user); err != nil {
		// Return validation errors
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, vErr := range validationErrors {
			errors[vErr.Field()] = vErr.Tag() // e.g., "required", "email", etc.
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	_, err := services.RegisterUser(user.Username, user.Password, user.Email)
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

func getUserDetails(c *gin.Context) {
	idStr := c.Param("id")
	userDetails, err := services.GetUserDetails(idStr) // Get userDetails from service layer
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, userDetails)
}

func getAllTodos(c *gin.Context) {
	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Assert userID to be a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asserting userID to string"})
		return
	}

	// Convert the string userID to a primitive.ObjectID
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID format"})
		return
	}

	todos, _ := services.GetTodos(objUserID) // Get todos from service layer
	c.JSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {
	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Assert userID to be a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asserting userID to string"})
		return
	}

	// Convert the string userID to a primitive.ObjectID
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID format"})
		return
	}

	idStr := c.Param("id")
	todos, err2 := services.GetTodoByID(idStr, objUserID) // Get todos from service layer
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	c.JSON(http.StatusOK, todos)
}
func createTodo(c *gin.Context) {
	var newTodo struct {
		Title string `bson:"title" json:"title" validate:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	// Validate the user struct
	if err := validate.Struct(&newTodo); err != nil {
		// Return validation errors
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, vErr := range validationErrors {
			errors[vErr.Field()] = vErr.Tag() // e.g., "required", "email", etc.
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Assert userID to be a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asserting userID to string"})
		return
	}

	// Convert the string userID to a primitive.ObjectID
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID format"})
		return
	}

	todoToAdd := models.Todo{Title: newTodo.Title}
	todoToAdd.ID = primitive.NewObjectID()
	todoToAdd.CreatedAt = time.Now()
	todoToAdd.UpdatedAt = time.Now()
	todoToAdd.UserID = objUserID

	result, err := services.AddTodo(todoToAdd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func updateTodo(c *gin.Context) {
	var newTodo models.TodoUpdate
	idStr := c.Param("id")
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	newTodo.UpdatedAt = time.Now()
	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Assert userID to be a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asserting userID to string"})
		return
	}

	// Convert the string userID to a primitive.ObjectID
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID format"})
		return
	}

	result, err := services.UpdateTodo(idStr, objUserID, newTodo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteTodo(c *gin.Context) {
	idStr := c.Param("id")

	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Assert userID to be a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error asserting userID to string"})
		return
	}

	// Convert the string userID to a primitive.ObjectID
	objUserID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID format"})
		return
	}

	_, err2 := services.DeleteTodo(idStr, objUserID)
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
