package user

//go:generate mockgen -source=usersvc.go -destination=../../../gen/mocks/mock_user/usersvc.go -self_package=../pkg/user

type UserData struct {
	Username      string
	EncryptedPass string
}

type UserSVC interface {
	PasswordEncrypt(pass string) (encryptedPass string, err error)
	AuthUser(username string, pass string) (bool, error)
	CreateUser(username string, pass string) error
}
