package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"userCrudApp/mocks"
	"userCrudApp/models"
)

func NewServer(userService *mocks.MockUserService) *gin.Engine {
	userController := NewUserController(userService)
	server := gin.Default()
	router := server.Group("/v1")
	userController.RegisterRoutes(router)
	return server
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user models.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser models.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}

func requireBodyMatchUsers(t *testing.T, body *bytes.Buffer, users []*models.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUsers []*models.User
	err = json.Unmarshal(data, &gotUsers)
	require.NoError(t, err)
	require.Equal(t, users, gotUsers)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
