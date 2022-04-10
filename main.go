package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"userCrudApp/config"
	"userCrudApp/controllers"
	"userCrudApp/services"
)

var (
	server         *gin.Engine
	userService    services.UserService
	userController controllers.UserController
	ctx            context.Context
	userCollection *mongo.Collection
	mongoClient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()

	mongoClient = config.ConnectDB()
	userCollection = config.GetCollection(mongoClient, "users")
	userService = services.NewUserServiceImpl(userCollection, ctx)
	userController = controllers.NewUserController(userService)
	server = gin.Default()
}
func main() {
	router := server.Group("/v1")
	userController.RegisterRoutes(router)
	log.Fatal(server.Run(":8000"))
}
