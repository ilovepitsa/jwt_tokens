package v1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	v1 "github.com/ilovepitsa/jwt_tokens/internal/http/v1"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type UserServiceMock struct {
	manager tokens.TokenManager
}

func (s *UserServiceMock) SignIn(user_id uint32, userIp string) (*entity.Tokens, error) {
	if user_id == 1 {
		return nil, fmt.Errorf("wrong id")
	}
	if user_id == 2 {
		return nil, nil
	}

	refresh, _ := s.manager.NewRefreshToken()
	return &entity.Tokens{AccessToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.RlGQ_kNzMgv-N_z_zaguUOgfgid8xMpvzmuqCBj6McOPcbXOCeTSdxFNtbCkXULDjYhU1UWV3ZXg7wRXu6CgWw",
		RefreshToken: refresh,
	}, nil
}
func (s *UserServiceMock) Refresh(refresh_toker string, userIp string) (*entity.Tokens, error) {
	if refresh_toker == "false" {
		return nil, fmt.Errorf("some error")
	}
	return &entity.Tokens{AccessToken: "access_token", RefreshToken: "refresh"}, nil
}
func (s *UserServiceMock) CreateUser() (*entity.User, error) {
	return nil, nil
}
func TestUserHandler(t *testing.T) {

	generator := func(b []byte) []byte {
		newB := make([]byte, len(b))
		var k uint8 = 0
		for i := range b {
			newB[i] = k
			k++
		}

		return newB
	}
	manager, err := tokens.NewManager([]byte("secret-key"), generator)
	if err != nil {
		t.Fatalf("cant create token manager")
	}
	testCases := []struct {
		name    string
		url     string
		service UserServiceMock
		want    string
		method  string
		test    func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string
	}{
		// SIGN IN METHOD
		{
			name: "Test wrong method sign",
			url:  "/auth/sign-in",

			service: UserServiceMock{manager: manager},
			method:  http.MethodPost,
			want:    "wrong method",
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
		{
			name: "Test sign-in without params",
			url:  "/auth/sign-in",

			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    `{"err": "add user_id input"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
		{
			name: "Test sign-in wrong id",
			url:  "/auth/sign-in?user_id=asdf",

			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    `{"err": "user id must be integer"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
		{
			name: "Test cant sign in ",
			url:  "/auth/sign-in?user_id=1",

			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    `{"err": "cant sign-in"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
		{
			name:    "Test sign-in answer",
			url:     "/auth/sign-in?user_id=3",
			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    `{"accessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.RlGQ_kNzMgv-N_z_zaguUOgfgid8xMpvzmuqCBj6McOPcbXOCeTSdxFNtbCkXULDjYhU1UWV3ZXg7wRXu6CgWw","refreshToken":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},

		// REFRESH METHOD
		{
			name:    "Test refresh wrong method",
			url:     "/auth/refresh",
			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    `{ "err" : "wrong method"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.Refresh(rr, r)
				return rr.Body.String()
			},
		},
		{
			name:    "Test refresh without params",
			url:     "/auth/refresh",
			method:  http.MethodPost,
			service: UserServiceMock{manager: manager},
			want:    `{"err":"need token"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.Refresh(rr, r)
				return rr.Body.String()
			},
		},
		{
			name:    "Test refresh cant refresh",
			url:     "/auth/refresh?token=false",
			method:  http.MethodPost,
			service: UserServiceMock{manager: manager},
			want:    `{"err": "cant refresh tokens"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.Refresh(rr, r)
				return rr.Body.String()
			},
		},
		{
			name:    "Test refresh answer",
			url:     "/auth/refresh?token=true",
			method:  http.MethodPost,
			service: UserServiceMock{manager: manager},
			want:    `{"accessToken":"access_token","refreshToken":"refresh"}`,
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.Refresh(rr, r)
				return rr.Body.String()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uh := v1.NewUserHandler(&tc.service)
			request, _ := http.NewRequest(tc.method, tc.url, nil)
			responce := httptest.NewRecorder()
			got := tc.test(uh, responce, request)
			got = strings.TrimRight(got, "\n")
			// t.Log(fmt.Printf("want: %v  got: %v", tc.want, got))
			if tc.want != got {
				// t.Fatalf("want: %v  got: %v", []byte(tc.want), []byte(got))
				t.Fatalf("want: %v  got: %v", tc.want, got)
			}
		})
	}

}
