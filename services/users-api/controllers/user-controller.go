package controllers

import (
	"fmt"
	"net/http"
	"users-api/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	// Implementar
	fmt.Printf("todo: implementar controllers ")
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	// Implementar
	ctx.JSON(http.StatusOK, gin.H{"message": "Get user by ID"})
}

func (c *UserController) Login(ctx *gin.Context) {
	// Implementar
	ctx.JSON(http.StatusOK, gin.H{"message": "Login"})
}
