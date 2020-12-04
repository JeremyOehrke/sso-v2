package session

import "time"

//go:generate mockgen -source=sessionsvc.go -destination=../../../gen/mocks/mock_session/sessionsvc.go -self_package=../pkg/sessionhandlers

const (
	MAX_SESSION_DURATION = 3600 * time.Second
)

type SessionData struct {
	Id          string            `json:"id"`
	Username    string            `json:"username"`
	SessionVars map[string]string `json:"sessionVars"`
}

type SessionSVC interface {
	GetSessionById(id string) (*SessionData, error)
	CreateSession(username string, sessionBody map[string]string) (sessionId string, err error)
	DestroySession(id string) error
	SetSessionBodyById(id string, body map[string]string) error
}

type SessionError string

func (e SessionError) Error() string { return string(e) }

const SessionNotFoundError = SessionError("session not found")
