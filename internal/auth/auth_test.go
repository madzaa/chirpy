package auth

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestCheckPasswordHash(t *testing.T) {
	getHash, _ := HashPassword("kokomoko")

	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "Validate hash for password", args: args{
			password: "kokomoko",
			hash:     getHash,
		}, want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPasswordHash(tt.args.password, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPasswordHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	id, err := uuid.Parse("2e4a206a-dc7d-412d-9da6-76ddfc981668")
	if err != nil {
		return
	}
	getJwt, err := MakeJWT(id, "i am the one who tests", 5*time.Minute)
	if err != nil {
		return
	}
	type args struct {
		tokenString string
		tokenSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{name: "Validate jwt", args: args{
			tokenString: getJwt,
			tokenSecret: "i am the one who tests",
		}, want: id, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.args.tokenString, tt.args.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
