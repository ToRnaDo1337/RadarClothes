package main

import (
	"net/http"

	"github.com/ToRnaDo1337/RadarClothes/login"
	"github.com/ToRnaDo1337/RadarClothes/registration"
	"github.com/ToRnaDo1337/RadarClothes/shop"
)

func main() {
	http.HandleFunc("/register", registration.Register)
	http.HandleFunc("/login", login.Login)
	http.HandleFunc("/addingItems", shop.HandleIndex)
	http.HandleFunc("/RadarClothes", shop.HandleItems)
	http.ListenAndServe(":8080", nil)
}
