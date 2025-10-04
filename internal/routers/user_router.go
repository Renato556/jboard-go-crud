package routers

import (
	"jboard-go-crud/internal/controllers"
	"net/http"
)

func NewUsersController(userHandler *controllers.UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users", userHandler.CreateUser)
	mux.HandleFunc("GET /v1/users", userHandler.GetUserHandler)
	mux.HandleFunc("PUT /v1/users", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /v1/users", userHandler.DeleteUser)
	return mux
}
