package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"userCrudApp/mocks"
	"userCrudApp/models"
)

func TestGetUser(t *testing.T) {
	user := models.User{
		Name: "Aaditya",
		Age:  21,
		Address: models.Address{
			City:    "Vapi",
			State:   "Gujarat",
			Pincode: 396191,
		},
	}

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(userService *mocks.MockUserService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "Success",
			username: user.Name,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					GetUser(&user.Name).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:     "InternalServerError",
			username: user.Name,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					GetUser(&user.Name).
					Times(1).
					Return(nil, mongo.ErrNoDocuments)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)

			server := NewServer(userService)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/users/getUser/%s", user.Name)
			request := httptest.NewRequest(http.MethodGet, url, nil)

			server.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateUser(t *testing.T) {
	user := models.User{
		Name: "Aaditya",
		Age:  21,
		Address: models.Address{
			City:    "Vapi",
			State:   "Gujarat",
			Pincode: 396191,
		},
	}

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(userService *mocks.MockUserService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					CreateUser(&user).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					CreateUser(gomock.Any()).
					Times(1).
					Return(errors.New(""))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					CreateUser(user).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)

			server := NewServer(userService)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/users/create")
			var request *http.Request
			if tc.name == "BadRequest" {
				request = httptest.NewRequest(http.MethodPost, url, nil)
			} else {
				request = httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			}

			server.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAll(t *testing.T) {
	userList := []*models.User{
		{
			Name: "Aaditya",
			Age:  21,
			Address: models.Address{
				City:    "Vapi",
				State:   "Gujarat",
				Pincode: 396191,
			},
		},
		{
			Name: "Akash",
			Age:  23,
			Address: models.Address{
				City:    "Surat",
				State:   "Gujarat",
				Pincode: 396600,
			},
		},
	}

	testCases := []struct {
		name          string
		buildStubs    func(userService *mocks.MockUserService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					GetAll().
					Times(1).
					Return(userList, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUsers(t, recorder.Body, userList)
			},
		},
		{
			name: "InternalServerError",
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					GetAll().
					Times(1).
					Return(nil, errors.New(""))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)

			server := NewServer(userService)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/users/getAll")
			var request *http.Request

			request = httptest.NewRequest(http.MethodGet, url, nil)

			server.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestUpdateUser(t *testing.T) {
	user := models.User{
		Name: "Aaditya",
		Age:  21,
		Address: models.Address{
			City:    "Vapi",
			State:   "Gujarat",
			Pincode: 396191,
		},
	}

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(userService *mocks.MockUserService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					UpdateUser(&user).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					UpdateUser(gomock.Any()).
					Times(1).
					Return(errors.New(""))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: user,
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					UpdateUser(user).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)

			server := NewServer(userService)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/users/update/%s", user.Name)
			var request *http.Request
			if tc.name == "BadRequest" {
				request = httptest.NewRequest(http.MethodPut, url, nil)
			} else {
				request = httptest.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			}

			server.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	username := "Aaditya"
	testCases := []struct {
		name          string
		buildStubs    func(userService *mocks.MockUserService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					DeleteUser(&username).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			buildStubs: func(userService *mocks.MockUserService) {
				userService.EXPECT().
					DeleteUser(&username).
					Times(1).
					Return(mongo.ErrNoDocuments)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)

			server := NewServer(userService)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/users/delete/%s", username)
			request := httptest.NewRequest(http.MethodDelete, url, nil)

			server.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}
