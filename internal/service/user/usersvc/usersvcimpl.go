package usersvc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sso-v2/internal/datasource"
	"sso-v2/internal/service/user"
)

type UserSVCImpl struct {
	ds datasource.Datasource
}

func NewUserSvc(ds datasource.Datasource) user.UserSVC {
	return &UserSVCImpl{ds: ds}
}

func (svc *UserSVCImpl) EncryptPassword(pass string) (encryptedPass string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		log.Print("error generating bcrypt hash: " + err.Error())
		return "", nil
	}

	//Encoded to handle weird binary issues when moving in and out of Redis
	encodedPass := base64.StdEncoding.EncodeToString(bytes)

	return encodedPass, nil
}

func (svc *UserSVCImpl) AuthUser(username string, pass string) (bool, error) {

	return false, nil
}

func (svc *UserSVCImpl) CreateUser(username string, encryptedPass string) error {
	foundUser, err := svc.ds.GetKey(generateUserKey(username))
	if err != nil {
		log.Print("error checking for existing username")
		return err
	}
	if foundUser != "" {
		return errors.New("username already taken")
	}

	userData := user.UserData{
		Username:      username,
		EncryptedPass: encryptedPass,
	}

	rawUser, err := json.Marshal(userData)
	if err != nil {
		log.Printf("error marshaling user data: %v", err.Error())
		return err
	}

	err = svc.ds.SetKey(generateUserKey(username), string(rawUser), 0)
	if err != nil {
		log.Printf("error writing user to datastore: %v", err.Error())
		return err
	}

	return nil
}

func generateUserKey(username string) string {
	return "user_" + username
}
