package v1

import (
	"net/http"

	"github.com/ilovepitsa/jwt_tokens/internal/service"
)

type UserHandlerInterface interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	userHandler UserHandlerInterface
}

func NewRouter(handler *http.ServeMux, services service.Services) {
	userHandler := NewUserHandler(services.UserService)
	router := &Router{
		userHandler: userHandler,
	}

	handler.HandleFunc("/auth/sign-in", router.userHandler.SignIn)
	handler.HandleFunc("/auth/refresh", router.userHandler.Refresh)
	handler.HandleFunc("/create", router.userHandler.CreateUser)
}
