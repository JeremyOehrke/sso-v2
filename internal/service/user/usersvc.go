package user

//go:generate mockgen -source=usersvc.go -destination=../../../gen/mocks/mock_user/usersvc.go -self_package=../pkg/userhandlers

type UserData struct {
	Username   string
	HashedPass string
}

type UserSVC interface {
	EncryptPassword(pass string) (encryptedPass string, err error)
	AuthUser(username string, pass string) (bool, error)
	CreateUser(username string, pass string) error
}

//Mapped Errors
type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

const NotFound = NotFoundError("user not found")
