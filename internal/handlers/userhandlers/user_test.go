package userhandlers

import (
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"sso-v2/gen/mocks/mock_user"
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
		method           string
		route            string
		requestBody      string
		username         string
		password         string
		expectSvcCall    bool
		expectedResponse expectedResponse
	}{
		{
			method:        "POST",
			route:         "/v1/users",
			name:          "missing everything",
			requestBody:   `{}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			method:        "POST",
			route:         "/v1/users",
			name:          "missing username",
			requestBody:   `{"password":"asdf"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			method:        "POST",
			route:         "/v1/users",
			name:          "missing password",
			requestBody:   `{"username":"joehrke"}`,
			expectSvcCall: false,
			expectedResponse: expectedResponse{
				statusCode: 400,
				body:       `{"message":"missing username and/or password"}`,
			},
		},
		{
			method:        "POST",
			route:         "/v1/users",
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

			ctrl := gomock.NewController(t)
			userSvc := mock_user.NewMockUserSVC(ctrl)
			if tt.expectSvcCall {
				userSvc.EXPECT().CreateUser(tt.username, tt.password)
			}

			router := apitest.BuildTestRouter(tt.method, tt.route, CreateUserHandler(userSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.route, strings.NewReader(tt.requestBody))
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
