package controller

import (
	"fmt"

	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiago-prog/novels-api/internal/model"
	"github.com/tiago-prog/novels-api/internal/usecase"
	"github.com/tiago-prog/novels-api/internal/utils"
)

type userController struct {
	userUseCase usecase.UserUsecase
	jwtSecret   string
}

func NewUserController(usecase usecase.UserUsecase) userController {
	return userController{
		userUseCase: usecase,
		jwtSecret:   os.Getenv("JWT_SECRET"),
	}
}

func (uc *userController) Register(ctx *gin.Context) {

	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if !utils.IsValidEmail(user.Email) {
		ctx.JSON(400, gin.H{"error": "Invalid email"})
		return
	}

	if !utils.IsValidUsername(user.Username) {
		ctx.JSON(400, gin.H{"error": "Invalid username"})
		return
	}

	id, err := uc.userUseCase.Register(user)
	if err != nil {
		fmt.Println("Error registering user:", err)
		ctx.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	// Create a payload with registration information
	payload := struct {
		Message string `json:"message"`
		UserID  int    `json:"user_id"`
	}{
		Message: "User registered successfully",
		UserID:  id,
	}

	ctx.JSON(200, payload)
}

func (uc *userController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user, err := uc.userUseCase.Login(loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to log in"})
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(user.ID), user.Role, []byte(uc.jwtSecret))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// Create a payload with login information
	payload := struct {
		Message string     `json:"message"`
		Token   string     `json:"token"`
		User    model.User `json:"user"`
	}{
		Message: "Login successful",
		Token:   token,
		User:    user,
	}

	ctx.JSON(200, payload)
}

func (uc *userController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	fmt.Println("Fetching user by email:", email)
	if email == "" {
		ctx.JSON(400, gin.H{"error": "Email is required"})
		return
	}

	user, err := uc.userUseCase.GetUserByEmail(email)
	if err != nil {
		fmt.Println("Error fetching user by email:", err)
		ctx.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Create a payload with user information
	payload := struct {
		Message string     `json:"message"`
		User    model.User `json:"user"`
	}{
		Message: "User fetched successfully",
		User:    user,
	}

	ctx.JSON(200, payload)
}

func (uc *userController) SuspendUser(ctx *gin.Context) {
	executorID := ctx.Param("executor_id")
	targetID := ctx.Param("target_id")
	fmt.Println("Suspending user:", executorID, targetID)
	if executorID == "" || targetID == "" {
		ctx.JSON(400, gin.H{"error": "Executor ID and target ID are required"})
		return
	}

	executorIDInt, err := strconv.Atoi(executorID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid executor ID"})
		return
	}
	targetIDInt, err := strconv.Atoi(targetID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := uc.userUseCase.SuspendUser(executorIDInt, targetIDInt); err != nil {
		fmt.Println("Error suspending user:", err)
		ctx.JSON(500, gin.H{"error": "Failed to suspend user"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User suspended successfully"})
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userUseCase.GetAllUsers()
	if err != nil {
		fmt.Println("Error fetching all users:", err)
		ctx.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Create a payload with user list
	payload := struct {
		Message string       `json:"message"`
		Users   []model.User `json:"users"`
	}{
		Message: "Users fetched successfully",
		Users:   users,
	}

	ctx.JSON(200, payload)
}
