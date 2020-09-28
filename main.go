package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hamidteimouri/go-jwt-restful/middlewares"

	"github.com/hamidteimouri/go-jwt-restful/controllers"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.Use(middlewares.JwtAuth)

	router.HandleFunc("/api/users/register", middlewares.SetJsonMiddleware(controllers.CreateUser)).Methods("POST")
	router.HandleFunc("/api/users/login", middlewares.SetJsonMiddleware(controllers.SignInUser)).Methods("POST")

	err := http.ListenAndServe(getPort(), router)

	if err != nil {
		fmt.Println(err)
	}
}

func getPort() string {
	if os.Getenv("PORT") != "" {
		return os.Getenv("PORT")
	}
	return ":8000"
}
