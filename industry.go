package main

import (
	"encoding/json"
	"net/http"

	// Import the cors package

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Industry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImgURL   string `json:"imgurl"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

func getIndustries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name, imgurl, title, subtitle, content FROM industries")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var industries []Industry
	for rows.Next() {
		var industry Industry
		err := rows.Scan(&industry.ID, &industry.Name, &industry.ImgURL, &industry.Title, &industry.Subtitle, &industry.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		industries = append(industries, industry)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(industries)
}

func createIndustry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newIndustry Industry
	err := json.NewDecoder(r.Body).Decode(&newIndustry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newIndustry.Name == "" {
		http.Error(w, "Industry name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO industries (name, imgurl, title, subtitle, content)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Execute the insert query and get the new ID
	err = db.QueryRow(query, newIndustry.Name, newIndustry.ImgURL, newIndustry.Title, newIndustry.Subtitle, newIndustry.Content).Scan(&newIndustry.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Industry already exists", http.StatusConflict)
				return

			}
		}
	}
	json.NewEncoder(w).Encode(newIndustry)
}

func updateIndustry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedIndustry Industry
	err := json.NewDecoder(r.Body).Decode(&updatedIndustry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE industries SET name = $1, imgurl = $2, title = $3, subtitle = $4, content = $5 WHERE id = $6;`
	result, err := db.Exec(query, updatedIndustry.Name, updatedIndustry.ImgURL, updatedIndustry.Title, updatedIndustry.Subtitle, updatedIndustry.Content, updatedIndustry.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check the affected rows", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Industry is not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updatedIndustry)
}

func deleteIndustry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	IndustryID := r.URL.Query().Get("id")

	query := `DELETE FROM industries WHERE id = $1;`
	result, err := db.Exec(query, IndustryID)

	if err != nil {
		http.Error(w, "Failed to delete Industry", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check the affected rows", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Industry is not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Industry deleted successfully"})
}
