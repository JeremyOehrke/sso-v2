package userhandlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"sso-v2/gen/mocks/mock_session"
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
	type userSvcAuthResponse struct {
		authed bool
		err    error
	}
	tests := []struct {
		name                    string
		requestBody             string
		username                string
		password                string
		expectUserSvcCall       bool
		expectedResponse        expectedResponse
		userSvcAuthResponse     userSvcAuthResponse
		expectSessionSvcCall    bool
		expectedSessionIdHeader string
		expectedSessionSvcError error
	}{
		{
			name:                 "missing everything",
			requestBody:          `{}`,
			expectUserSvcCall:    false,
			expectSessionSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:                 "missing username",
			requestBody:          `{"password":"asdf"}`,
			expectUserSvcCall:    false,
			expectSessionSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:                 "missing password",
			requestBody:          `{"username":"joehrke"}`,
			expectUserSvcCall:    false,
			expectSessionSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			name:              "Auth Failed user found",
			expectUserSvcCall: true,
			username:          "joehrke",
			password:          "asdf",
			requestBody:       `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":false}`,
			},
			userSvcAuthResponse: userSvcAuthResponse{
				authed: false,
				err:    nil,
			},
			expectSessionSvcCall: false,
		},
		{
			name:              "Auth Failed user not found",
			expectUserSvcCall: true,
			username:          "joehrke",
			password:          "asdf",
			requestBody:       `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":false}`,
			},
			userSvcAuthResponse: userSvcAuthResponse{
				authed: false,
				err:    user.NotFound,
			},
			expectSessionSvcCall: false,
		},
		{
			name:              "Auth Failed odd error",
			expectUserSvcCall: true,
			username:          "joehrke",
			password:          "asdf",
			requestBody:       `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 500,
				body:       `{"message":"error authorizing user"}`,
			},
			userSvcAuthResponse: userSvcAuthResponse{
				authed: false,
				err:    errors.New("some weird error"),
			},
			expectSessionSvcCall: false,
		},
		{
			name:              "Auth Success",
			expectUserSvcCall: true,
			username:          "joehrke",
			password:          "asdf",
			requestBody:       `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 200,
				body:       `{"authOk":true}`,
			},
			userSvcAuthResponse: userSvcAuthResponse{
				authed: true,
				err:    nil,
			},
			expectSessionSvcCall:    true,
			expectedSessionIdHeader: "asdf-1235",
			expectedSessionSvcError: nil,
		},
		{
			name:              "Auth Success With Session Create Error",
			expectUserSvcCall: true,
			username:          "joehrke",
			password:          "asdf",
			requestBody:       `{"username":"joehrke","password":"asdf"}`,
			expectedResponse: expectedResponse{
				statusCode: 500,
				body:       `{"message":"error creating session"}`,
			},
			userSvcAuthResponse: userSvcAuthResponse{
				authed: true,
				err:    nil,
			},
			expectSessionSvcCall:    true,
			expectedSessionIdHeader: "",
			expectedSessionSvcError: errors.New("some weird error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "POST"
			url := "/v1/user/doAuth"

			ctrl := gomock.NewController(t)
			userSvc := mock_user.NewMockUserSVC(ctrl)
			if tt.expectUserSvcCall {
				userSvc.EXPECT().AuthUser(tt.username, tt.password).Return(tt.userSvcAuthResponse.authed, tt.userSvcAuthResponse.err)
			}

			sessionSvc := mock_session.NewMockSessionSVC(ctrl)
			if tt.expectSessionSvcCall {
				sessionSvc.EXPECT().SetSession(tt.username, gomock.Any()).Return(tt.expectedSessionIdHeader, tt.expectedSessionSvcError)
			}

			router := apitest.BuildTestRouter(method, url, AuthUserHandler(userSvc, sessionSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, url, strings.NewReader(tt.requestBody))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedResponse.statusCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedResponse.statusCode)
			}
			if w.Header().Get(SessionIdHeader) != tt.expectedSessionIdHeader {
				t.Errorf("Unexpected session id header -- got %v, wanted %v", w.Header().Get(SessionIdHeader), tt.expectedSessionIdHeader)
			}
			if strings.TrimSuffix(w.Body.String(), "\n") != tt.expectedResponse.body {
				t.Errorf("Unexpected body -- got: %v, wanted: %v", w.Body.String(), tt.expectedResponse.body)
			}
		})
	}
}
