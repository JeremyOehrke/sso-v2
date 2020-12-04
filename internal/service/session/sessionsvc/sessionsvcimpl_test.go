package sessionsvc

import (
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"sso-v2/gen/mocks/mock_datasource"
	"sso-v2/internal/service/session"
	"testing"
)

func Test_generateSessionKey(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Base",
			args: args{id: "12345"},
			want: "sess_12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateSessionKey(tt.args.id); got != tt.want {
				t.Errorf("generateSessionKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionSVCImpl_CreateSession(t *testing.T) {
	type args struct {
		username    string
		sessionBody map[string]string
	}
	tests := []struct {
		name              string
		args              args
		minExpectedLength int
		wantErr           bool
		err               error
	}{
		{
			name: "HappyPath",
			args: args{
				username: "joehrke",
				sessionBody: map[string]string{
					"test_key": "val",
				},
			},
			minExpectedLength: 36,
			wantErr:           false,
			err:               nil,
		},
		{
			name: "Redis_Error",
			args: args{
				username: "joehrke",
				sessionBody: map[string]string{
					"test_key": "val",
				},
			},
			minExpectedLength: 0,
			wantErr:           true,
			err:               errors.New("test redis error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)
			ds.EXPECT().SetKey(gomock.Any(), gomock.Any(), session.MAX_SESSION_DURATION).Return(tt.err)

			svc := &SessionSVCImpl{
				ds: ds,
			}
			gotSessionId, err := svc.CreateSession(tt.args.username, tt.args.sessionBody)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotSessionId) != tt.minExpectedLength {
				t.Errorf("CreateSession() gotSessionId length: %v, want %v", len(gotSessionId), tt.minExpectedLength)
			}
			ctrl.Finish()
		})
	}
}

func TestSessionSVCImpl_DestroySession(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "Happy_Path",
			args:    args{id: "12345"},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "Redis_Error",
			args:    args{id: "12345"},
			wantErr: true,
			err:     errors.New("test Redis error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)
			ds.EXPECT().DelKey(generateSessionKey(tt.args.id)).Return(tt.err)

			svc := &SessionSVCImpl{
				ds: ds,
			}
			if err := svc.DestroySession(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DestroySession() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctrl.Finish()
		})
	}
}

func TestSessionSVCImpl_GetSessionById_EmptySession(t *testing.T) {
	ctrl := gomock.NewController(t)
	ds := mock_datasource.NewMockDatasource(ctrl)
	ds.EXPECT().GetKey(generateSessionKey("12345")).Return("", nil)

	svc := &SessionSVCImpl{
		ds: ds,
	}
	sess, err := svc.GetSessionById("12345")
	if err != nil && err != session.SessionNotFoundError {
		t.Errorf("erronious error")
		return
	}
	if sess != nil {
		t.Errorf("returned sessionhandlers when one doesn't exist")
		return
	}
	ctrl.Finish()
}

func TestSessionSVCImpl_GetSessionById_GetKeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ds := mock_datasource.NewMockDatasource(ctrl)
	ds.EXPECT().GetKey(generateSessionKey("12345")).Return("", errors.New("test error"))

	svc := &SessionSVCImpl{
		ds: ds,
	}
	sess, err := svc.GetSessionById("12345")
	if err == nil {
		t.Errorf("should have gotten error")
		return
	}
	if sess != nil {
		t.Errorf("returned sessionhandlers when one doesn't exist")
		return
	}
	ctrl.Finish()
}

func TestSessionSVCImpl_GetSessionById_SessionFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	ds := mock_datasource.NewMockDatasource(ctrl)
	ds.EXPECT().GetKey(generateSessionKey("12345")).Return(`{"key":"12345", "username":"joehrke", "sessionVars":{"test":"val"}}`, nil)

	//This expect makes sure key expiration is reset in Redis
	ds.EXPECT().SetKey(generateSessionKey("12345"), `{"key":"12345", "username":"joehrke", "sessionVars":{"test":"val"}}`, session.MAX_SESSION_DURATION)

	svc := &SessionSVCImpl{
		ds: ds,
	}
	sess, err := svc.GetSessionById("12345")
	if err != nil {
		t.Errorf("erronious error")
		return
	}
	expectedSession := &session.SessionData{
		Id:          "12345",
		Username:    "joehrke",
		SessionVars: map[string]string{"test": "val"},
	}
	if !reflect.DeepEqual(sess, expectedSession) {
		t.Error("expected sessionhandlers didn't match")
		return
	}
	ctrl.Finish()
}

func TestSessionSVCImpl_SetSessionBodyById(t *testing.T) {
	type getSessionRequest struct {
		expected      bool
		key           string
		responseBody  string
		responseError error
	}
	type setSessionRequest struct {
		expected      bool
		key           string
		responseError error
	}
	type args struct {
		id   string
		body map[string]string
	}
	tests := []struct {
		name              string
		args              args
		getSessionRequest getSessionRequest
		setSessionRequest setSessionRequest
		wantErr           bool
		respError         error
	}{
		{
			name: "session not found",
			args: args{
				id: "asdf-1234",
			},
			getSessionRequest: getSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseBody:  "",
				responseError: session.SessionNotFoundError,
			},
			setSessionRequest: setSessionRequest{
				expected: false,
			},
			wantErr:   true,
			respError: session.SessionNotFoundError,
		},
		{
			name: "get session redis error",
			args: args{
				id: "asdf-1234",
			},
			getSessionRequest: getSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseBody:  "",
				responseError: errors.New("some weird error"),
			},
			setSessionRequest: setSessionRequest{
				expected: false,
			},
			wantErr:   true,
			respError: errors.New("some weird error"),
		},
		{
			name: "set session redis error",
			args: args{
				id: "asdf-1234",
			},
			getSessionRequest: getSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseBody:  `{"id":"asdf-1234","username":"joehrke","sessionVars":{"test":"val"}}`,
				responseError: nil,
			},
			setSessionRequest: setSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseError: errors.New("some weird Redis error"),
			},
			wantErr:   true,
			respError: errors.New("some weird error"),
		},
		{
			name: "all ok",
			args: args{
				id: "asdf-1234",
			},
			getSessionRequest: getSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseBody:  `{"id":"asdf-1234","username":"joehrke","sessionVars":{"test":"val"}}`,
				responseError: nil,
			},
			setSessionRequest: setSessionRequest{
				expected:      true,
				key:           "sess_asdf-1234",
				responseError: nil,
			},
			wantErr:   false,
			respError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)

			if tt.getSessionRequest.expected {
				ds.EXPECT().GetKey(tt.getSessionRequest.key).Return(tt.getSessionRequest.responseBody, tt.getSessionRequest.responseError)
			}

			if tt.setSessionRequest.expected {
				ds.EXPECT().SetKey(gomock.Any(), gomock.Any(), session.MAX_SESSION_DURATION).Return(tt.respError)
			}

			svc := &SessionSVCImpl{
				ds: ds,
			}
			err := svc.SetSessionBodyById(tt.args.id, tt.args.body)

			ctrl.Finish()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSessionBodyById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
