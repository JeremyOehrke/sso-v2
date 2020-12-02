package usersvc

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
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
			gotEncryptedPass, err := svc.PasswordEncrypt(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decodedHash, err := base64.StdEncoding.DecodeString(gotEncryptedPass)
			if err != nil {
				t.Error("mis-encoded password hash")
			}

			if bcrypt.CompareHashAndPassword(decodedHash, []byte(tt.args.pass)) != nil {
				t.Errorf("PasswordEncrypt() want hashes to match")
			}
		})
	}
}
