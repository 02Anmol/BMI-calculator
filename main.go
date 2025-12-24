package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

// --- Constants and File Path ---
// const dataFile = "users_data.json"

// --- Data Model ---
type User struct {
	Name     string  `json:"name"`
	WeightKg float64 `json:"weight_kg"`
	HeightM  float64 `json:"height_m"`
	BMI      float64 `json:"bmi"`
	Category string  `json:"category"`
}

// ViewModel is used to pass data to the HTML template.
type ViewModel struct {
	Users   []User
	Message string // For displaying success/error messages
}

// Global variable to hold all user records in memory.
var users []User

// Global template variable. Must be parsed once at startup.
var tpl *template.Template

// --- Backend (File Operations) ---

// loadUserData attempts to read and unmarshal the JSON data from the file.
func loadUserData() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			users = []User{}
			log.Printf("Note: %s not found. Starting with an empty user list.", dataFile)
			return
		}
		log.Fatalf("Error reading data file: %v", err)
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}
	log.Printf("Loaded %d user records from %s.", len(users), dataFile)
}

// saveUserData marshals the current 'users' slice and writes it back to the file.
func saveUserData() error {
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling user data: %w", err)
	}

	err = os.WriteFile(dataFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}
	return nil
}

// --- BMI Calculation Functions ---

// calculateBMI computes the Body Mass Index.
func calculateBMI(weightKg float64, heightM float64) float64 {
	if heightM <= 0 {
		return 0.0
	}

	return weightKg / (heightM * heightM)
}

// getBMICategory returns a categorical interpretation of the calculated BMI.
func getBMICategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi >= 18.5 && bmi <= 24.9:
		return "Normal Weight"
	case bmi >= 25.0 && bmi <= 29.9:
		return "Overweight"
	case bmi >= 30.0:
		return "Obesity"
	default:
		return "Cannot interpret"
	}
}

// --- HTTP Handlers ---

// indexHandler displays the main page with the form and the data table.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Prepare the data to be passed to the template
	data := ViewModel{
		Users: users, // Pass the current list of users
	}

	// 2. Execute the template
	err := tpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// calculateHandler processes the form submission, calculates BMI, saves data, and redirects.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Extract and validate input
	name := r.FormValue("name")
	weightStr := r.FormValue("weight")
	heightStr := r.FormValue("height")

	weightKg, errW := strconv.ParseFloat(weightStr, 64)
	heightM, errH := strconv.ParseFloat(heightStr, 64)

	if errW != nil || errH != nil || weightKg <= 0 || heightM <= 0 {
		http.Error(w, "Invalid input. Please enter valid positive numbers for weight and height.", http.StatusBadRequest)
		return
	}

	// 3. Calculate BMI and Category
	bmi := calculateBMI(weightKg, heightM)
	category := getBMICategory(bmi)

	// 4. Create new User record
	newUser := User{
		Name:     name,
		WeightKg: weightKg,
		HeightM:  heightM,
		BMI:      bmi,
		Category: category,
	}

	// 5. Store data
	users = append(users, newUser)

	// 6. Save all data to the file (backend)
	if err := saveUserData(); err != nil {
		log.Printf("Failed to save data: %v", err)
		// Still redirect, but log the error
	}

	// 7. Redirect back to the index page
	http.Redirect(w, r, "/?status=success", http.StatusSeeOther)
}

func main() {
	// 1. Initialize: Load data and parse templates
	loadUserData()
	var err error
	// Parses all files in the templates folder that end with .html
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	// 2. Define HTTP routes (Endpoints)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This handles the status message after a successful POST request
		if r.URL.Query().Get("status") == "success" {
			data := ViewModel{
				Users:   users,
				Message: fmt.Sprintf("Success! %s's BMI (%.2f) calculated and saved.", users[len(users)-1].Name, users[len(users)-1].BMI),
			}
			if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
				http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		indexHandler(w, r)
	})
	http.HandleFunc("/calculate", calculateHandler)

	// 3. Start the server
	port := ":8080"
	log.Printf("Starting web server on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}