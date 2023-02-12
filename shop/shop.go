package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Item struct {
	Name  string
	Price float64
}

var items []Item

func ItemPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("item.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			return
		}
		item := Item{
			Name:  r.FormValue("name"),
			Price: price,
		}
		items = append(items, item)
		http.Redirect(w, r, "/item", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/item", ItemPageHandler)
	http.HandleFunc("/item/add", AddItemHandler)
	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
