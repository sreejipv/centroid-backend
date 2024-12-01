package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// Import the cors package

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Award struct {
	ID       string `json:"id"`
	Awardurl string `json:"awardurl"`
	Type     string `json:"type"`
	Name     string `json:"name"`
}
type Client struct {
	ID        string `json:"id"`
	Clienturl string `json:"clienturl"`
	Name      string `json:"name"`
}

func createAward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newAward Award
	err := json.NewDecoder(r.Body).Decode(&newAward)

	if (err) != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newAward.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO awards ( name, awardurl, type)
	VALUES ($1,$2,$3)
	RETURNING id`

	err = db.QueryRow(query, newAward.Name, newAward.Awardurl, newAward.Type).Scan(newAward.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Award already exists", http.StatusConflict) // HTTP 409 Conflict
				return
			}
		}
	}

	json.NewEncoder(w).Encode(newAward)
}
func fetchAward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing award ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, award, type FROM awards WHERE id=$1`
	var award Award

	err := db.QueryRow(query, id).Scan(&award.ID, &award.Name, &award.Awardurl, &award.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "client not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(award); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getAwards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name, awardurl, type FROM awards")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var awards []Award
	for rows.Next() {
		var award Award
		err := rows.Scan(&award.ID, &award.Name, &award.Awardurl, &award.Type)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		awards = append(awards, award)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(awards)
}
func updateAward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedAward Award
	err := json.NewDecoder(r.Body).Decode(&updatedAward)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedAward.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `UPDATE awards SET name = $1, awardurl = $2, type = $3 WHERE id = $4;`
	result, err := db.Exec(query, updatedAward.Name, updatedAward.Awardurl, updatedAward.Type, updatedAward.ID)

	if err != nil {
		http.Error(w, "Failed to update award", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "award not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedAward)
}

func deleteAward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	awardID := r.URL.Query().Get("id")

	query := `DELETE FROM awards WHERE id = $1;`
	result, err := db.Exec(query, awardID)

	if err != nil {
		http.Error(w, "Failed to delete award", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to check affected rows ", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "award is not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "award deleted successfully"})
}

func createClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newClient Client
	err := json.NewDecoder(r.Body).Decode(&newClient)

	if (err) != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newClient.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO clients ( name, clienturl)
	VALUES ($1,$2)
	RETURNING id`

	err = db.QueryRow(query, newClient.Name, newClient.Clienturl).Scan(newClient.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Client already exists", http.StatusConflict) // HTTP 409 Conflict
				return
			}
		}
	}

	json.NewEncoder(w).Encode(newClient)
}

func getClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name, clienturl FROM clients")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var clients []Client
	for rows.Next() {
		var client Client
		err := rows.Scan(&client.ID, &client.Name, &client.Clienturl)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(clients)
}

func fetchClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, clienturl FROM clients WHERE id=$1`
	var client Client

	err := db.QueryRow(query, id).Scan(&client.ID, &client.Name, &client.Clienturl)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "client not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func updateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedClient Client
	err := json.NewDecoder(r.Body).Decode(&updatedClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedClient.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `UPDATE clients SET name = $1, clienturl = $2 WHERE id = $3;`
	result, err := db.Exec(query, updatedClient.Name, updatedClient.Clienturl, updatedClient.ID)

	if err != nil {
		http.Error(w, "Failed to update Client", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedClient)
}

func deleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientID := r.URL.Query().Get("id")

	query := `DELETE FROM clients WHERE id = $1;`
	result, err := db.Exec(query, clientID)

	if err != nil {
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to check affected rows ", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "client is not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "client deleted successfully"})
}
