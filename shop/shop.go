package shop

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"strings"
)

type Item struct {
	Title       string
	Description string
	Photo       string
}

// func main() {

// 	http.ListenAndServe(":8080", nil)
// }

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "item.html")
		return
	}

	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20) // limit 32 MB

		title := r.FormValue("title")
		description := r.FormValue("description")
		photo := r.FormValue("photoURL")

		f, err := os.OpenFile("./uploads/"+photo, os.O_WRONLY|os.O_CREATE, 0666)

		defer f.Close()

		// Save the form data to a file
		data := fmt.Sprintf("%s,%s,%s\n", title, description, photo)
		if err := os.MkdirAll("./data", 0755); err != nil {
			http.Error(w, "Error creating data directory", http.StatusBadRequest)
			return
		}
		f, err = os.OpenFile("./data/items.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			http.Error(w, "Error opening data file", http.StatusBadRequest)
			return
		}
		defer f.Close()
		if _, err := io.WriteString(f, data); err != nil {
			http.Error(w, "Error saving data", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/RadarClothes", http.StatusSeeOther)
	}
}

func HandleItems(w http.ResponseWriter, r *http.Request) {
	// Read the saved data from the file
	data, err := ioutil.ReadFile("./data/items.txt")
	if err != nil {
		http.Error(w, "Error reading data", http.StatusBadRequest)
		return
	}

	// Split the data into individual items
	items := strings.Split(string(data), "\n")

	// Set the number of items per page
	itemsPerPage := 3

	// Calculate the total number of pages
	numItems := len(items)
	numPages := numItems / itemsPerPage
	if numItems%itemsPerPage != 0 {
		numPages++
	}

	// Get the current page number from the URL query parameters
	currentPage := 1
	if r.URL.Query().Get("page") != "" {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err == nil && page > 0 && page <= numPages {
			currentPage = page
		}
	}

	// Get the search query from the URL query parameters
	searchQuery := r.URL.Query().Get("search")

	// Filter the items based on the search query
	filteredItems := make([]string, 0)
	for _, item := range items {
		if searchQuery != "" {
			if !strings.Contains(strings.ToLower(item), strings.ToLower(searchQuery)) {
				continue
			}
		}
		if item == "" {
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	// Get the slice of items for the current page
	start := (currentPage - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > len(filteredItems) {
		end = len(filteredItems)
	}
	pageItems := filteredItems[start:end]

	// Generate the HTML for the page
	var html strings.Builder
	html.WriteString("<html><body>")

	// Add the main.html template
	mainTmpl, err := template.ParseFiles("main.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusBadRequest)
		return
	}

	// Define a struct to hold the data that will be used to fill the template
	type ItemData struct {
		Title       string
		Description string
		Photo       string
	}

	// Create a slice of ItemData structs to hold the data for each item
	itemData := make([]ItemData, len(pageItems))

	// Fill the slice with the data for each item
	for i, item := range pageItems {
		fields := strings.Split(item, ",")
		title := fields[0]
		description := fields[1]
		photo := fields[2]

		itemData[i] = ItemData{title, description, photo}
	}

	// Execute the template, passing in the slice of ItemData structs and pagination variables
	err = mainTmpl.Execute(&html, itemData)

	if err != nil {
		http.Error(w, "Error executing template", http.StatusBadRequest)
		return
	}

	html.WriteString("</body></html>")

	// Write the HTML to the response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html.String()))
}
