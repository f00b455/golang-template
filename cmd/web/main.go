package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/f00b455/golang-template/internal/config"
	"github.com/f00b455/golang-template/internal/handlers"
	"github.com/f00b455/golang-template/pkg/shared"
)

type PageData struct {
	Title     string
	Headlines []shared.RssHeadline
	UpdatedAt string
	Error     string
}

var templates *template.Template

func main() {
	// Load config (currently unused, but may be needed for future settings)
	_ = config.Load()

	// Parse templates
	funcMap := template.FuncMap{
		"formatDate": formatDate,
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	// Set up routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/headlines", headlinesAPIHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Web server starting on port %s", port)
	log.Printf("Visit http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start web server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch headlines from API
	headlines, err := fetchHeadlines()

	data := PageData{
		Title:     "SPIEGEL Headlines",
		Headlines: headlines,
		UpdatedAt: time.Now().Format("15:04:05"),
	}

	if err != nil {
		data.Error = "Unable to fetch headlines"
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func headlinesAPIHandler(w http.ResponseWriter, r *http.Request) {
	headlines, err := fetchHeadlines()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unable to fetch headlines"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"headlines": headlines,
		"updatedAt": time.Now().Format(time.RFC3339),
	})
}

func fetchHeadlines() ([]shared.RssHeadline, error) {
	// Fetch from the API server
	apiURL := "http://localhost:3002/api/rss/spiegel/top5"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response handlers.HeadlinesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Headlines, nil
}

func formatDate(dateStr string) string {
	// Parse the date
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}

	// Convert to Berlin timezone
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.Local
	}

	return t.In(loc).Format("02.01.2006 15:04")
}