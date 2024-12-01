package main

import (
	"encoding/json"
	"net/http"

	// Import the cors package

	_ "github.com/lib/pq"
)

type Contact struct {
	ID          string `json:"id"`
	FullName    string `json:"fullname"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CompanyName string `json:"companyname"`
	EmailID     string `json:"emailid"`
	Phone       string `json:"phone"`
	Requirement string `json:"requirement"`
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, fullname, city, country, companyname, emailid, phone, requirement FROM contact")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		err := rows.Scan(&contact.ID, &contact.FullName, &contact.City, &contact.CompanyName, &contact.Country, &contact.EmailID, &contact.Phone, &contact.Requirement)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(contacts)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newContact Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newContact.FullName == "" {
		http.Error(w, "Contact name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO contact (fullname, city, country, companyname, emailid, phone, requirement)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	// Execute the insert query and get the new ID
	err = db.QueryRow(query, newContact.FullName, newContact.City, newContact.CompanyName, newContact.Country, newContact.EmailID, newContact.Phone, newContact.Requirement).Scan(&newContact.ID)

	if err != nil {
		http.Error(w, "Failed to create contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newContact)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ContactID := r.URL.Query().Get("id")

	query := `DELETE FROM contact WHERE id = $1;`
	result, err := db.Exec(query, ContactID)

	if err != nil {
		http.Error(w, "Failed to delete Contact", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check the affected rows", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Contact is not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Contact deleted successfully"})
}
