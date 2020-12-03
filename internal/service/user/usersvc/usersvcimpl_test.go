package usersvc

import (
	"encoding/base64"
	"errors"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"sso-v2/gen/mocks/mock_datasource"
	"testing"
)

func TestUserSVCImpl_PasswordEncrypt(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "HappyPath",
			args:    args{pass: "abc123"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &UserSVCImpl{
				ds: nil,
			}
			gotEncryptedPass, err := svc.EncryptPassword(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decodedHash, err := base64.StdEncoding.DecodeString(gotEncryptedPass)
			if err != nil {
				t.Error("password encoding issue")
			}

			if bcrypt.CompareHashAndPassword(decodedHash, []byte(tt.args.pass)) != nil {
				t.Errorf("EncryptPassword() want hashes to match")
			}
		})
	}
}

func TestUserSVCImpl_CreateUser_NameTaken(t *testing.T) {
	type args struct {
		username string
		pass     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		dsErr   error
	}{
		{
			name: "Redis_Error",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			wantErr: true,
			dsErr:   errors.New("test redis error"),
		},
		{
			name: "Name_In_Use",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			wantErr: true,
			dsErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)
			ds.EXPECT().GetKey(generateUserKey(tt.args.username)).Return(`{"username":"joehrke", "encryptedPass":"sdgsdfsdfsdf"`, tt.dsErr)

			svc := &UserSVCImpl{
				ds: ds,
			}

			if err := svc.CreateUser(tt.args.username, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestUserSVCImpl_CreateUser(t *testing.T) {
	type args struct {
		username string
		pass     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		dsErr   error
	}{
		{
			name: "HappyPath",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			wantErr: false,
			dsErr:   nil,
		},
		{
			name: "Redis_Error",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			wantErr: true,
			dsErr:   errors.New("test Redis error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)
			ds.EXPECT().GetKey(generateUserKey(tt.args.username)).Return("", nil)
			ds.EXPECT().SetKey(generateUserKey(tt.args.username), gomock.Any(), 0).Return(tt.dsErr)

			svc := &UserSVCImpl{
				ds: ds,
			}

			if err := svc.CreateUser(tt.args.username, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestUserSVCImpl_AuthUser(t *testing.T) {
	type args struct {
		username string
		pass     string
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		userFound string
		dsError   error
		wantErr   bool
	}{
		{
			name: "User Not Found",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			userFound: "",
			dsError:   nil,
			wantErr:   true,
			want:      false,
		},
		{
			name: "Redis Error",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			userFound: "",
			dsError:   errors.New("test Redis error"),
			wantErr:   true,
			want:      false,
		},
		{
			name: "Password Mismatch",
			args: args{
				username: "joehrke",
				pass:     "abc1234",
			},
			userFound: `{"username":"joehrke", "hashedPass":"JDJhJDE0JHBmRnZYWWtvaXViWEFIdm9ra2lEVi5PbzhwRXppQXNBOHJzNXVyYXpIZ2NYM1BjczlkMy91"}`,
			dsError:   nil,
			wantErr:   false,
			want:      false,
		},
		{
			name: "Password Match",
			args: args{
				username: "joehrke",
				pass:     "abc123",
			},
			userFound: `{"username":"joehrke", "hashedPass":"JDJhJDE0JHBmRnZYWWtvaXViWEFIdm9ra2lEVi5PbzhwRXppQXNBOHJzNXVyYXpIZ2NYM1BjczlkMy91"}`,
			dsError:   nil,
			wantErr:   false,
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ds := mock_datasource.NewMockDatasource(ctrl)
			ds.EXPECT().GetKey(generateUserKey(tt.args.username)).Return(tt.userFound, tt.dsError)

			svc := &UserSVCImpl{
				ds: ds,
			}

			got, err := svc.AuthUser(tt.args.username, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
