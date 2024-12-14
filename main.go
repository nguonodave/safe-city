package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"safecity/vars"
)

// IncidentReport represents a single report
type IncidentReport struct {
	IncidentType   string `json:"incidentType"`
	DateTime       string `json:"dateTime"`
	Location       string `json:"location"`
	Description    string `json:"description"`
	AdditionalInfo string `json:"additionalInfo"`
	ContactEmail   string `json:"contactEmail"`
}

// Reports storage (in-memory)
var (
	reports  []IncidentReport
	mutex    sync.Mutex
	dataFile = "reports.json"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars.AllTemplates.ExecuteTemplate(w, "home.html", nil)
}

func emergencyContactsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars.AllTemplates.ExecuteTemplate(w, "emergency.html", nil)
}

// loadReports loads existing reports from the JSON file
func loadReports() error {
	// Check if file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil // we will return nil because the file does not exist yet
	}

	// Read the file
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	// If file is empty, initialize empty reports slice
	if len(data) == 0 {
		reports = make([]IncidentReport, 0)
		return nil
	}

	// Unmarshal JSON into reports slice
	return json.Unmarshal(data, &reports)
}

// saveReports saves the reports to the JSON file
func saveReports() error {
    data, err := json.MarshalIndent(reports, "", "    ")
    if err != nil {
        return err
    }
    return os.WriteFile(dataFile, data, 0644)
}

func submitReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars.AllTemplates.ExecuteTemplate(w, "report.html", nil)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Create an IncidentReport object from form data
	report := IncidentReport{
		IncidentType:   r.FormValue("incidentType"),
		DateTime:       r.FormValue("dateTime"),
		Location:       r.FormValue("location"),
		Description:    r.FormValue("description"),
		AdditionalInfo: r.FormValue("additionalInfo"),
		ContactEmail:   r.FormValue("contactEmail"),
	}

	// Safely append to the slice and save to file using mutex
    mutex.Lock()
    reports = append(reports, report)
    err = saveReports()
    mutex.Unlock()

	// Redirect to success page
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

// successReport renders the success page to the success path
func successReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars.AllTemplates.ExecuteTemplate(w, "success.html", nil)
}

func getReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Return all reports as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func main() {
	// Initialize reports by loading from file
    if err := loadReports(); err != nil {
        log.Printf("Error loading reports: %v", err)
    }

	vars.AllTemplates, _ = vars.AllTemplates.ParseGlob(vars.TemplatesDir + "*.html")

	http.HandleFunc("/", homePage) // Render home page
	http.HandleFunc("/emergency-contacts", emergencyContactsPage) // Render emergency contacts page
	http.HandleFunc("/submit-report", submitReport) // Render report submission
	http.HandleFunc("/success", successReport)      // Render success report submission
	http.HandleFunc("/report", handleReport)        // Handle form submissions
	http.HandleFunc("/reports", getReports)         // Endpoint to fetch reports

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
