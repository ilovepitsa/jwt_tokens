package v1_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	v1 "github.com/ilovepitsa/jwt_tokens/internal/http/v1"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type UserServiceMock struct {
	manager tokens.TokenManager
}

func (s *UserServiceMock) SignIn(r *http.Request) (*entity.Tokens, error) {
	return nil, nil
}
func (s *UserServiceMock) Refresh(r *http.Request) (*entity.Tokens, error) {
	return nil, nil
}
func (s *UserServiceMock) CreateUser(r *http.Request) (*entity.User, error) {
	return nil, nil
}
func TestUserHandler(t *testing.T) {

	type answer struct {
		token entity.Tokens
	}
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
		input   io.Reader
		service UserServiceMock
		want    string
		method  string
		test    func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string
	}{
		{
			name:    "Test wrong method sign",
			input:   strings.NewReader(""),
			service: UserServiceMock{manager: manager},
			method:  http.MethodPost,
			want:    "wrong method",
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
		{
			name:    "Test wrong method sign",
			input:   strings.NewReader("12345"),
			method:  http.MethodGet,
			service: UserServiceMock{manager: manager},
			want:    "wrong method",
			test: func(uh *v1.UserHandler, rr *httptest.ResponseRecorder, r *http.Request) string {
				uh.SignIn(rr, r)
				return rr.Body.String()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uh := v1.NewUserHandler(&tc.service)
			request, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", input)
			responce := httptest.NewRecorder()
			got := tc.test(uh, responce, request)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatal("want: ", tc.want, " got: ", got)
			}
		})
	}

}
