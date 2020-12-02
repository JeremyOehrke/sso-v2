package usersvc

import (
	"encoding/base64"
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

func (svc *UserSVCImpl) PasswordEncrypt(pass string) (encryptedPass string, err error) {
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

func (svc *UserSVCImpl) CreateUser(username string, pass string) error {

	return nil
}
