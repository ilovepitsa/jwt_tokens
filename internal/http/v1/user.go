package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	if r.Method != http.MethodGet {
		log.Println("wrong method request sign-in")
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	if !r.URL.Query().Has("user_id") {
		log.Println("Wrong input")
		http.Error(w, `{"err": "add user_id input"}`, http.StatusBadRequest)
		return
	}
	id_str := r.URL.Query().Get("user_id")
	user_id, err := strconv.ParseUint(id_str, 10, 32)
	if err != nil {
		log.Println("Cant conver user_id")
		http.Error(w, `{"err": "user id must be integer"}`, http.StatusBadRequest)
		return
	}
	res, err := uh.userService.SignIn(uint32(user_id))
	if err != nil {
		log.Println("cant sign-in")
		http.Error(w, `{"err": "cant sign-in"}`, http.StatusInternalServerError)
		return
	}

	ans, err := json.Marshal(res)
	if err != nil {
		log.Println("cant marshal tokens")

		http.Error(w, `{"err": "cant marshal tokens"}`, http.StatusInternalServerError)
		return
	}
	w.Write(ans)

}

func (uh *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{ "err" : "wrong method"}`, http.StatusMethodNotAllowed)
		return
	}
	if !r.URL.Query().Has("token") {
		log.Println("need token")
		http.Error(w, `{"err":"need token"}`, http.StatusBadRequest)
		return
	}

	refresh := r.URL.Query().Get("token")
	res, err := uh.userService.Refresh(refresh)
	if err != nil {
		log.Println("cant refresh ", err)
		http.Error(w, `{"err": "cant refresh tokens"}`, http.StatusInternalServerError)
		return
	}

	ans, err := json.Marshal(res)
	if err != nil {
		log.Println("cant marshal tokens")

		http.Error(w, `{"err": "cant marshal tokens"}`, http.StatusInternalServerError)
		return
	}
	w.Write(ans)

}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{ "err" : "wrong method"}`, http.StatusMethodNotAllowed)
		return
	}
	user, err := uh.userService.CreateUser()
	if err != nil {
		log.Println(err)
		http.Error(w, `{ "err" : "cant create user"}`, http.StatusInternalServerError)
		return
	}
	ans := fmt.Sprintf(`{ "id": %d}`, user.Id)
	w.Write([]byte(ans))

}
