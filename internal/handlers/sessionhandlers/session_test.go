package sessionhandlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"sso-v2/gen/mocks/mock_session"
	"sso-v2/internal/service/session"
	"sso-v2/internal/test/apitest"
	"strings"
	"testing"
)

func TestGetSessionDataHandler(t *testing.T) {
	type getSessionByIdRequest struct {
		expected    bool
		requestedId string
		sessionData *session.SessionData
		err         error
	}
	type expectedHttpResponse struct {
		statusCode int
		jsonBody   string
	}
	type requestData struct {
		route    string
		jsonBody string
	}
	tests := []struct {
		name                  string
		route                 string
		method                string
		requestData           requestData
		getSessionByIdRequest getSessionByIdRequest
		expectedHttpResponse  expectedHttpResponse
	}{
		{
			name:   "no session id provided",
			route:  "/session/:sessionId",
			method: "GET",
			requestData: requestData{
				route:    "/session/",
				jsonBody: "",
			},
			getSessionByIdRequest: getSessionByIdRequest{
				expected:    false,
				requestedId: "",
				sessionData: nil,
				err:         nil,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 404,
				jsonBody:   "404 page not found",
			},
		},
		{
			name:   "session id not found",
			route:  "/session/:sessionId",
			method: "GET",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: "",
			},
			getSessionByIdRequest: getSessionByIdRequest{
				expected:    true,
				requestedId: "asdf-1234",
				sessionData: nil,
				err:         session.SessionNotFoundError,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 404,
				jsonBody:   "",
			},
		},
		{
			name:   "session id lookup failure",
			route:  "/session/:sessionId",
			method: "GET",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: "",
			},
			getSessionByIdRequest: getSessionByIdRequest{
				expected:    true,
				requestedId: "asdf-1234",
				sessionData: nil,
				err:         errors.New("some weird error"),
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 500,
				jsonBody:   `{"message":"error locating session"}`,
			},
		},
		{
			name:   "session id found",
			route:  "/session/:sessionId",
			method: "GET",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: "",
			},
			getSessionByIdRequest: getSessionByIdRequest{
				expected:    true,
				requestedId: "asdf-1234",
				sessionData: &session.SessionData{
					Id:       "asdf-1234",
					Username: "joehrke",
					SessionVars: map[string]string{
						"test": "val",
					},
				},
				err: nil,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 200,
				jsonBody:   `{"id":"asdf-1234","username":"joehrke","sessionVars":{"test":"val"}}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sessionSvc := mock_session.NewMockSessionSVC(ctrl)
			if tt.getSessionByIdRequest.expected {
				sessionSvc.EXPECT().GetSessionById(tt.getSessionByIdRequest.requestedId).Return(tt.getSessionByIdRequest.sessionData, tt.getSessionByIdRequest.err)
			}

			router := apitest.BuildTestRouter(tt.method, tt.route, GetSessionDataHandler(sessionSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.requestData.route, strings.NewReader(tt.requestData.jsonBody))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedHttpResponse.statusCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedHttpResponse.statusCode)
			}
			if strings.TrimSuffix(w.Body.String(), "\n") != tt.expectedHttpResponse.jsonBody {
				t.Errorf("Unexpected body -- got: %v, wanted: %v", strings.TrimSuffix(w.Body.String(), "\n"), tt.expectedHttpResponse.jsonBody)
			}
		})
	}
}

func TestSetSessionDataHandler(t *testing.T) {
	type setSessionRequest struct {
		expected    bool
		sessionId   string
		sessionBody interface{}
		err         error
	}
	type expectedHttpResponse struct {
		statusCode int
		jsonBody   string
	}
	type requestData struct {
		route    string
		jsonBody string
	}
	tests := []struct {
		name                 string
		route                string
		method               string
		requestData          requestData
		setSessionRequest    setSessionRequest
		expectedHttpResponse expectedHttpResponse
	}{
		{
			name:   "no session ID passed",
			route:  "/session/:sessionId",
			method: "PUT",
			requestData: requestData{
				route:    "/session/",
				jsonBody: "",
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 404,
				jsonBody:   "404 page not found",
			},
		},
		{
			name:   "missing body parts",
			route:  "/session/:sessionId",
			method: "PUT",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: "{}",
			},
			setSessionRequest: setSessionRequest{
				expected: false,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 400,
				jsonBody:   `{"message":"missing body"}`,
			},
		},
		{
			name:   "session not found",
			route:  "/session/:sessionId",
			method: "PUT",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: `{"sessionVars":{"test":"val"}}`,
			},
			setSessionRequest: setSessionRequest{
				expected:  true,
				sessionId: "asdf-1234",
				sessionBody: map[string]string{
					"test": "val",
				},
				err: session.SessionNotFoundError,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 404,
				jsonBody:   "",
			},
		},
		{
			name:   "OK",
			route:  "/session/:sessionId",
			method: "PUT",
			requestData: requestData{
				route:    "/session/asdf-1234",
				jsonBody: `{"sessionVars":{"test":"val"}}`,
			},
			setSessionRequest: setSessionRequest{
				expected:  true,
				sessionId: "asdf-1234",
				sessionBody: map[string]string{
					"test": "val",
				},
				err: nil,
			},
			expectedHttpResponse: expectedHttpResponse{
				statusCode: 204,
				jsonBody:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sessionSvc := mock_session.NewMockSessionSVC(ctrl)
			if tt.setSessionRequest.expected {
				sessionSvc.EXPECT().SetSessionBodyById(tt.setSessionRequest.sessionId, tt.setSessionRequest.sessionBody).Return(tt.setSessionRequest.err)
			}

			router := apitest.BuildTestRouter(tt.method, tt.route, SetSessionDataHandler(sessionSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.requestData.route, strings.NewReader(tt.requestData.jsonBody))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedHttpResponse.statusCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedHttpResponse.statusCode)
			}
			if strings.TrimSuffix(w.Body.String(), "\n") != tt.expectedHttpResponse.jsonBody {
				t.Errorf("Unexpected body -- got: %v, wanted: %v", strings.TrimSuffix(w.Body.String(), "\n"), tt.expectedHttpResponse.jsonBody)
			}
		})
	}
}

func TestDestroySessionHandler(t *testing.T) {
	type destroySessionRequest struct {
		expected  bool
		sessionId string
		err       error
	}
	tests := []struct {
		name                  string
		route                 string
		requestRoute          string
		destroySessionRequest destroySessionRequest
		expectedResponseCode  int
	}{
		{
			name:         "session destroy error",
			route:        "/sessions/:sessionId",
			requestRoute: "/sessions/asdf-1234",
			destroySessionRequest: destroySessionRequest{
				expected:  true,
				sessionId: "asdf-1234",
				err:       errors.New("some weird errror"),
			},
			expectedResponseCode: http.StatusInternalServerError,
		},
		{
			name:         "session not found error",
			route:        "/sessions/:sessionId",
			requestRoute: "/sessions/asdf-1234",
			destroySessionRequest: destroySessionRequest{
				expected:  true,
				sessionId: "asdf-1234",
				err:       session.SessionNotFoundError,
			},
			expectedResponseCode: http.StatusOK, //This is OK because it effectively means the session was already gone
		},
		{
			name:         "session destroyed",
			route:        "/sessions/:sessionId",
			requestRoute: "/sessions/asdf-1234",
			destroySessionRequest: destroySessionRequest{
				expected:  true,
				sessionId: "asdf-1234",
				err:       nil,
			},
			expectedResponseCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sessionSvc := mock_session.NewMockSessionSVC(ctrl)
			if tt.destroySessionRequest.expected {
				sessionSvc.EXPECT().DestroySession(tt.destroySessionRequest.sessionId).Return(tt.destroySessionRequest.err)
			}

			router := apitest.BuildTestRouter("DELETE", tt.route, DestroySessionHandler(sessionSvc))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tt.requestRoute, strings.NewReader(""))
			router.ServeHTTP(w, req)

			ctrl.Finish()
			if w.Code != tt.expectedResponseCode {
				t.Errorf("Unexpected status code -- got: %v, wanted: %v", w.Code, tt.expectedResponseCode)
			}
		})
	}
}
