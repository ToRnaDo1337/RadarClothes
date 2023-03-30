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

type Rating struct {
	ItemID   int
	Stars    int
	Comments string
}

var ratings []Rating

type ItemData struct {
	Title       string
	Description string
	Photo       string
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "item.html")
		return
	}

	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)

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

	data, err := ioutil.ReadFile("./data/items.txt")
	if err != nil {
		http.Error(w, "Error reading data", http.StatusBadRequest)
		return
	}

	// Split the data into individual items
	items := strings.Split(string(data), "\n")

	// Set the number of items per page
	itemsPerPage := 1000

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
	price, err := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
	if err != nil {
		price = 0
	}
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
		fields := strings.Split(item, ",")
		if len(fields) < 4 {
			continue
		}
		if price > 0 {
			itemPrice, err := strconv.ParseFloat(fields[3], 64)
			if err != nil || itemPrice > price {
				continue
			}
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

func HandleRateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		// Get the item ID from the form data
		itemID, err := strconv.Atoi(r.FormValue("itemID"))
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		// Get the rating from the form data
		stars, err := strconv.Atoi(r.FormValue("stars"))
		if err != nil {
			http.Error(w, "Invalid rating", http.StatusBadRequest)
			return
		}
		comments := r.FormValue("comments")

		// Add the new rating to the ratings slice
		newRating := Rating{ItemID: itemID, Stars: stars, Comments: comments}
		ratings = append(ratings, newRating)

		// Save the ratings to file
		f, err := os.OpenFile("./data/Rating.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			http.Error(w, "Error opening rating file", http.StatusBadRequest)
			return
		}
		defer f.Close()

		for _, rating := range ratings {
			data := fmt.Sprintf("%d,%d,%s\n", rating.ItemID, rating.Stars, rating.Comments)
			if _, err := io.WriteString(f, data); err != nil {
				http.Error(w, "Error saving rating data", http.StatusBadRequest)
				return
			}
		}

		http.Redirect(w, r, "/RadarClothes/items", http.StatusSeeOther)
	}
}

func RateHandler(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.FormValue("itemID")
	starsStr := r.FormValue("stars")
	comments := r.FormValue("comments")

	// Parse item ID and rating
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	stars, err := strconv.Atoi(starsStr)
	if err != nil || stars < 1 || stars > 5 {
		http.Error(w, "Invalid rating", http.StatusBadRequest)
		return
	}

	// Write rating to file
	f, err := os.OpenFile("Rating.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open rating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d %d %s\n", itemID, stars, comments)
	if err != nil {
		http.Error(w, "Failed to write rating to file", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getRating(itemID int) Rating {
	for _, rating := range ratings {
		if rating.ItemID == itemID {
			return rating
		}
	}
	return Rating{}
}
