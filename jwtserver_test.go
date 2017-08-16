package yljwt

import (
	"net/http"
	"reflect"
	"testing"
)

var jwt256 *JwtServer

func TestJwtServer_Token(t *testing.T) {
	type args struct {
		user  LoginUser
		store TokenStore
	}

	jwt256 = NewJwtServerHS256("appidceshi", "aaaaa", Header, 1000)
	tests := []struct {
		name            string
		j               JwtServer
		args            args
		wantTokenstring string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "测试h256",
			j:    *jwt256,
			args: args{
				user: UserInfo{LoginKind: "wechat", LoginID: "aaaaaa"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenstring, err := tt.j.Token(tt.args.user, tt.args.store)
			t.Log(err)
			t.Log(gotTokenstring)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("JwtServer.Token() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if gotTokenstring != tt.wantTokenstring {
			// 	t.Errorf("JwtServer.Token() = %v, want %v", gotTokenstring, tt.wantTokenstring)
			// }
		})
	}
}

func TestJwtServer_WriteToken(t *testing.T) {
	type args struct {
		w     http.ResponseWriter
		user  LoginUser
		store TokenStore
	}
	tests := []struct {
		name    string
		j       JwtServer
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.WriteToken(tt.args.w, tt.args.user, tt.args.store); (err != nil) != tt.wantErr {
				t.Errorf("JwtServer.WriteToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJwtServer_CheckToken(t *testing.T) {
	type args struct {
		r         *http.Request
		checkuser func(user UserInfo) bool
	}
	tests := []struct {
		name     string
		j        JwtServer
		args     args
		wantUser UserInfo
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := tt.j.CheckToken(tt.args.r, tt.args.checkuser)
			if (err != nil) != tt.wantErr {
				t.Errorf("JwtServer.CheckToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("JwtServer.CheckToken() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
