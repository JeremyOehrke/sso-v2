package usersvc

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sso-v2/internal/datasource"
	"sso-v2/internal/service/user"
)

const ()

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
	//encodedPass := base64.StdEncoding.EncodeToString(bytes)

	return string(bytes), nil
}

func (svc *UserSVCImpl) AuthUser(username string, pass string) (bool, error) {
	foundUser, err := svc.ds.GetKey(generateUserKey(username))
	if err != nil {
		log.Print("error checking for existing username")
		return false, err
	}
	if foundUser == "" {
		return false, user.NotFound
	}

	//Check passwords match
	userDat := &user.UserData{}
	err = json.Unmarshal([]byte(foundUser), userDat)
	if err != nil {
		log.Printf("error unmarshaling userhandlers data: %v", err.Error())
		return false, err
	}

	//hashPass, err := base64.StdEncoding.DecodeString(userDat.HashedPass)
	//if err != nil {
	//	log.Printf("error base64 decoding: %v", err.Error())
	//	return false, err
	//}

	err = bcrypt.CompareHashAndPassword([]byte(userDat.HashedPass), []byte(pass))
	//if our passwords mismatch, its not a failure, just rejected
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil { //This catches all other errors and returns the error
		log.Printf("password comparison error: %v", err.Error())
		return false, err
	}

	//If the check has made it this far, all's well
	return true, nil
}

func (svc *UserSVCImpl) CreateUser(username string, encryptedPass string) error {
	foundUser, err := svc.ds.GetKey(generateUserKey(username))
	if err != nil && err != datasource.KeyNotFound {
		log.Print("error checking for existing username")
		return err
	}
	if foundUser != "" {
		return errors.New("username already taken")
	}

	userData := user.UserData{
		Username:   username,
		HashedPass: encryptedPass,
	}

	rawUser, err := json.Marshal(userData)
	if err != nil {
		log.Printf("error marshaling userhandlers data: %v", err.Error())
		return err
	}

	err = svc.ds.SetKey(generateUserKey(username), string(rawUser), 0)
	if err != nil {
		log.Printf("error writing userhandlers to datastore: %v", err.Error())
		return err
	}

	return nil
}

func generateUserKey(username string) string {
	return "user_" + username
}
