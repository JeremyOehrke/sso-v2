package sessionsvc

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"sso-v2/internal/datasource"
	"sso-v2/internal/service/session"
)

type SessionSVCImpl struct {
	ds datasource.Datasource
}

func NewSessionSvc(ds datasource.Datasource) session.SessionSVC {
	return &SessionSVCImpl{
		ds: ds,
	}
}

func (svc *SessionSVCImpl) GetSessionById(id string) (*session.SessionData, error) {
	rawSess, err := svc.ds.GetKey(generateSessionKey(id))
	if err != nil {
		log.Print("error fetching sessionhandlers by key: " + err.Error())
		return nil, err
	}
	// if our key returns empty, no sessionhandlers exists
	if len(rawSess) == 0 {
		return nil, session.SessionNotFoundError
	}

	//if we have a sessionhandlers, unmarshal it and return
	sess := &session.SessionData{}
	err = json.Unmarshal([]byte(rawSess), sess)
	if err != nil {
		log.Print("error unmarshaling sessionhandlers data: " + err.Error())
		return nil, err
	}

	//bump sessionhandlers expiration in Redis
	err = svc.ds.SetKey(generateSessionKey(id), rawSess, session.MAX_SESSION_DURATION)
	if err != nil {
		log.Print("error resetting sessionhandlers timeout: " + err.Error())
		return nil, err //returns nil even if sessionhandlers found to ensure no strange behavior
	}

	return sess, nil
}

func (svc *SessionSVCImpl) CreateSession(username string, sessionBody map[string]string) (sessionId string, err error) {
	sessionId = generateSessionId()
	sess := session.SessionData{
		Id:          generateSessionKey(sessionId),
		Username:    username,
		SessionVars: sessionBody,
	}
	rawSess, err := json.Marshal(sess)
	if err != nil {
		log.Print("error marshaling sessionhandlers data: " + err.Error())
		return "", err
	}

	err = svc.ds.SetKey(generateSessionKey(sessionId), string(rawSess), session.MAX_SESSION_DURATION)
	if err != nil {
		log.Print("error writing key to store: " + err.Error())
		return "", err
	}

	return sessionId, nil
}

func (svc *SessionSVCImpl) DestroySession(id string) error {
	err := svc.ds.DelKey(generateSessionKey(id))
	if err != nil {
		log.Print("error deleting key from store: " + err.Error())
		return err
	}
	return nil
}

func (svc *SessionSVCImpl) SetSessionBodyById(id string, body map[string]string) error {
	rawSess, err := svc.ds.GetKey(generateSessionKey(id))
	if err != nil {
		log.Print("error fetching sessionhandlers by key: " + err.Error())
		return err
	}
	// if our key returns empty, no sessionhandlers exists
	if len(rawSess) == 0 {
		return session.SessionNotFoundError
	}

	//if we have a sessionhandlers, unmarshal it and return
	sess := &session.SessionData{}
	err = json.Unmarshal([]byte(rawSess), sess)
	if err != nil {
		log.Print("error unmarshaling sessionhandlers data: " + err.Error())
		return err
	}

	//Set the body and remarshal back into a string
	sess.SessionVars = body
	sessBytes, err := json.Marshal(sess)
	if err != nil {
		log.Print("error marshaling sessionhandlers data: " + err.Error())
		return err
	}

	//reset the key
	err = svc.ds.SetKey(generateSessionKey(id), string(sessBytes), session.MAX_SESSION_DURATION)
	if err != nil {
		log.Print("error setting key: " + err.Error())
		return err
	}

	return nil
}

func generateSessionId() string {
	return uuid.New().String()
}

func generateSessionKey(id string) string {
	return "sess_" + id
}
