package userhandlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"sso-v2/gen/mocks/mock_user"
	"sso-v2/internal/service/user"
	"sso-v2/internal/test/apitest"
	"strings"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	type expectedResponse struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name             string
		requestBody      string
		username         string
		password         string
		expectSvcCall    bool
		expectedResponse expectedResponse
	}{
		{
			name:          "missing everything",
			requestBody:   `{}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "missing username",
			requestBody:   `{"password":"asdf"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "missing password",
			requestBody:   `{"username":"joehrke"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "OK",
			expectSvcCall: true,
			username:      "joehrke",
			password:      "asdf",
			requestBody:   `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 201,
				body:       ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "POST"
			url := "/v1/user/doAuth"

			ctrl := gomock.NewController(t)
			userSvc := mock_user.NewMockUserSVC(ctrl)
			if tt.expectSvcCall {
				userSvc.EXPECT().CreateUser(tt.username, tt.password)
			}

			router := apitest.BuildTestRouter(method, url, CreateUserHandler(userSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, url, strings.NewReader(tt.requestBody))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedResponse.statusCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedResponse.statusCode)
			}
			if strings.TrimSuffix(w.Body.String(), "\n") != tt.expectedResponse.body {
				t.Errorf("Unexpected body -- got: %v, wanted: %v", w.Body.String(), tt.expectedResponse.body)
			}
		})
	}
}

func TestAuthUserHandler(t *testing.T) {
	type expectedResponse struct {
		statusCode int
		body       string
	}
	type svcAuthResponse struct {
		authed bool
		err    error
	}
	tests := []struct {
		name             string
		requestBody      string
		username         string
		password         string
		expectSvcCall    bool
		expectedResponse expectedResponse
		svcAuthResponse  svcAuthResponse
	}{
		{
			name:          "missing everything",
			requestBody:   `{}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "missing username",
			requestBody:   `{"password":"asdf"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "missing password",
			requestBody:   `{"username":"joehrke"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:          "Auth Failed user found",
			expectSvcCall: true,
			username:      "joehrke",
			password:      "asdf",
			requestBody:   `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":false}`,
			},
			svcAuthResponse: svcAuthResponse{
				authed: false,
				err:    nil,
			},
		},
		{
			name:          "Auth Failed user not found",
			expectSvcCall: true,
			username:      "joehrke",
			password:      "asdf",
			requestBody:   `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":false}`,
			},
			svcAuthResponse: svcAuthResponse{
				authed: false,
				err:    user.NotFound,
			},
		},
		{
			name:          "Auth Failed odd error",
			expectSvcCall: true,
			username:      "joehrke",
			password:      "asdf",
			requestBody:   `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 500,
				body:       `{"message":"error authorizing user"}`,
			},
			svcAuthResponse: svcAuthResponse{
				authed: false,
				err:    errors.New("some weird error"),
			},
		},
		{
			name:          "Auth Success",
			expectSvcCall: true,
			username:      "joehrke",
			password:      "asdf",
			requestBody:   `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":true}`,
			},
			svcAuthResponse: svcAuthResponse{
				authed: true,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "POST"
			url := "/v1/user/doAuth"

			ctrl := gomock.NewController(t)
			userSvc := mock_user.NewMockUserSVC(ctrl)
			if tt.expectSvcCall {
				userSvc.EXPECT().AuthUser(tt.username, tt.password).Return(tt.svcAuthResponse.authed, tt.svcAuthResponse.err)
			}

			router := apitest.BuildTestRouter(method, url, AuthUserHandler(userSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, url, strings.NewReader(tt.requestBody))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedResponse.statusCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedResponse.statusCode)
			}
			if strings.TrimSuffix(w.Body.String(), "\n") != tt.expectedResponse.body {
				t.Errorf("Unexpected body -- got: %v, wanted: %v", w.Body.String(), tt.expectedResponse.body)
			}
		})
	}
}
