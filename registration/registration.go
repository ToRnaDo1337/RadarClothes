package registration

import (
	"html/template"
	"net/http"
	"os"
)

type User struct {
	Username string
	Password string
	Email    string
}

var tmpl = template.Must(template.ParseFiles("register.html"))

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.Execute(w, nil)
		return
	}

	user := User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	file, err := os.OpenFile("users.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(user.Username + "," + user.Password + "," + user.Email + "\n"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
