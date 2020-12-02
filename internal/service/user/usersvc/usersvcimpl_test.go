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
				t.Error("mis-encoded password hash")
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
