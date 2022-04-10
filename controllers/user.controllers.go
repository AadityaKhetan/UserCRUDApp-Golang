package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userCrudApp/models"
	"userCrudApp/services"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{UserService: userService}
}

func (uc UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.CreateUser(&user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func (uc UserController) GetUser(c *gin.Context) {
	username := c.Param("name")
	user, err := uc.UserService.GetUser(&username)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc UserController) GetAll(c *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (uc UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func (uc UserController) DeleteUser(c *gin.Context) {
	username := c.Param("name")
	if err := uc.UserService.DeleteUser(&username); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (uc UserController) RegisterRoutes(rg *gin.RouterGroup) {
	userRouter := rg.Group("/users")
	userRouter.POST("/create", uc.CreateUser)
	userRouter.GET("/getUser/:name", uc.GetUser)
	userRouter.GET("/getAll", uc.GetAll)
	userRouter.PUT("/update/:name", uc.UpdateUser)
	userRouter.DELETE("/delete/:name", uc.DeleteUser)
}
