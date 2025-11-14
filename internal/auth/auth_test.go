package auth

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
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
	getJwt, err := MakeJWT(id, "i am the one who tests")
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Bearer Token",
			header:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			wantErr: false,
		},
		{
			name:    "Missing Authorization Header",
			header:  "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Missing Bearer Prefix",
			header:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Wrong Prefix (Basic Auth)",
			header:  "Basic dXNlcjpwYXNz",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Bearer with No Token",
			header:  "Bearer ",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Lowercase bearer",
			header:  "bearer token123",
			want:    "",
			wantErr: true, // Unless your implementation is case-insensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create header for each test case
			headers := http.Header{}
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetBearerToken(headers)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetBearerToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	// The function returns a non-empty string.
	//The string is 64 hex characters (since 32 bytes * 2).
	//The string is valid hex.
	//No error is returned.
	// dst := make([]byte, hex.DecodedLen(len(src)))
	//n, err := hex.Decode(dst, src)
	//if err != nil {
	//log.Fatal(err)

	token, err := MakeRefreshToken()

	if err != nil {
		t.Errorf("can't make refresh token %v", err)
		return
	}

	dst := make([]byte, hex.DecodedLen(len(token)))
	_, err = hex.Decode(dst, []byte(token))
	if err != nil {
		t.Errorf("can't  decode token %v", err)
		return
	}

	if len(token) != 64 {
		t.Errorf("string is not 64 hex character, it's %v", len(token))
		return
	}
}
