package main

import (
	"net/http"

	"github.com/ToRnaDo1337/RadarClothes/login"
	"github.com/ToRnaDo1337/RadarClothes/registration"
)

func main() {
	http.HandleFunc("/register", registration.Register)
	http.HandleFunc("/login", login.Login)
	http.ListenAndServe(":8080", nil)
}
