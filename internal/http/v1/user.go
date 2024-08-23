package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilovepitsa/jwt_tokens/internal/service"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Println("here request")
	if r.Method != http.MethodGet {
		log.Println("wrong method request sign-in")
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}
	res, err := uh.userService.SignIn(r)
	if err != nil {
		log.Println("cant sign-in")
		http.Error(w, "cant sign-in", http.StatusInternalServerError)
		return
	}

	ans, err := json.Marshal(res)
	if err != nil {
		log.Println("cant marshal tokens")
		http.Error(w, "cant marshal tokens", http.StatusInternalServerError)
		return
	}
	w.Write(ans)

}

func (uh *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
