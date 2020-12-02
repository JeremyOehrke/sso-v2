package session

//go:generate mockgen -source=sessionsvc.go -destination=../../../gen/mocks/mock_session/sessionsvc.go -self_package=../pkg/session

const (
	MAX_SESSION_DURATION = 3600
)

type SessionData struct {
	Id          string
	Username    string
	SessionVars map[string]string
}

type SessionSVC interface {
	GetSessionById(id string) (*SessionData, error)
	CreateSession(username string, sessionBody map[string]string) (sessionId string, err error)
	DestroySession(id string) error
}
