package main

import (
	"encoding/json"
	"log"
	"net/http"

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
var reports []IncidentReport

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

	// Create an IncidentReport object from form data
	report := IncidentReport{
		IncidentType:   r.FormValue("incidentType"),
		DateTime:       r.FormValue("dateTime"),
		Location:       r.FormValue("location"),
		Description:    r.FormValue("description"),
		AdditionalInfo: r.FormValue("additionalInfo"),
		ContactEmail:   r.FormValue("contactEmail"),
	}

	reports = append(reports, report)

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
	vars.AllTemplates, _ = vars.AllTemplates.ParseGlob(vars.TemplatesDir + "*.html")

	http.HandleFunc("/submit-report", submitReport) // Render report submission
	http.HandleFunc("/success", successReport)      // Render report submission
	http.HandleFunc("/report", handleReport)        // Handle form submissions
	http.HandleFunc("/reports", getReports)         // Endpoint to fetch reports

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
