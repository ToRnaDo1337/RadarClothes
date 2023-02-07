package login

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")

		file, err := os.Open("users.txt")
		if err != nil {
			fmt.Fprint(w, "Error opening users file")
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				continue
			}

			if parts[0] == username && parts[1] == password {
				fmt.Fprintf(w, "Welcome, %s!", username)
				return
			}
		}

		fmt.Fprint(w, "Invalid username or password")
	}
}
