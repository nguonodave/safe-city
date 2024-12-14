package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
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
	reports      []IncidentReport
	mutex        sync.Mutex
	dataFile     = "reports.json"
	AllTemplates *template.Template
	TemplatesDir = "templates/"
)

func handleError(w http.ResponseWriter, r *http.Request, statusCode int) {
	var templateName string
	switch statusCode {
	case http.StatusNotFound:
		templateName = "404.html"
	case http.StatusInternalServerError:
		templateName = "500.html"
	case http.StatusMethodNotAllowed:
		templateName = "405.html"
	case http.StatusBadRequest:
		templateName = "400.html"
	default:
		templateName = "error.html" // Generic error page
	}

	// Render error template based on status code
	w.WriteHeader(statusCode)
	err := AllTemplates.ExecuteTemplate(w, templateName, nil)
	if err != nil {
		log.Printf("Error rendering error page for status %d: %v", statusCode, err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	// Render the home page template
	AllTemplates.ExecuteTemplate(w, "home.html", nil)
}

func overviewPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	AllTemplates.ExecuteTemplate(w, "overview.html", nil)
}

func emergencyContactsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	AllTemplates.ExecuteTemplate(w, "emergency.html", nil)
}

func submitReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	AllTemplates.ExecuteTemplate(w, "report.html", nil)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handleError(w, r, http.StatusBadRequest)
		return
	}

	report := IncidentReport{
		IncidentType:   r.FormValue("incidentType"),
		DateTime:       r.FormValue("dateTime"),
		Location:       r.FormValue("location"),
		Description:    r.FormValue("description"),
		AdditionalInfo: r.FormValue("additionalInfo"),
		ContactEmail:   r.FormValue("contactEmail"),
	}

	mutex.Lock()
	reports = append(reports, report)
	saveErr := saveReports()
	mutex.Unlock()

	if saveErr != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

func successReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	AllTemplates.ExecuteTemplate(w, "success.html", nil)
}

func loadReports() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		reports = make([]IncidentReport, 0)
		return nil
	}

	return json.Unmarshal(data, &reports)
}

func saveReports() error {
	data, err := json.MarshalIndent(reports, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0o644)
}

func getReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func main() {
	if err := loadReports(); err != nil {
		log.Printf("Error loading reports: %v", err)
	}

	var err error
	AllTemplates, err = template.ParseGlob(TemplatesDir + "*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			homePage(w, r)
			return
		}
		handleError(w, r, http.StatusNotFound)
	})
	http.HandleFunc("/overview", overviewPage)
	http.HandleFunc("/emergency-contacts", emergencyContactsPage)
	http.HandleFunc("/submit-report", submitReport)
	http.HandleFunc("/success", successReport)
	http.HandleFunc("/report", handleReport)
	http.HandleFunc("/reports", getReports)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
