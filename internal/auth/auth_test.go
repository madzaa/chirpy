package auth

import (
	"testing"
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
		{name: "Hash password", args: args{
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
