package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// Import the cors package

	_ "github.com/lib/pq"
)

type Banner struct {
	ID            string `json:"id"`
	DeskBannerUrl string `json:"deskbannerurl"`
	MobBannerUrl  string `json:"mobbannerurl"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	CtaText       string `json:"ctatext"`
	CtaAction     string `json:"ctaaction"`
}

var banners []Banner

func getBanners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, deskbannerurl, mobbannerurl, title, subtitle, ctatext, ctaaction FROM banners")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var banners []Banner
	for rows.Next() {
		var banner Banner
		err := rows.Scan(&banner.ID, &banner.DeskBannerUrl, &banner.MobBannerUrl, &banner.Title, &banner.Subtitle, &banner.CtaText, &banner.CtaAction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(banners)
}

func fetchBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing banner ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, deskbannerurl, mobbannerurl, title, subtitle, ctatext, ctaaction FROM banners WHERE ID=$1`
	var banner Banner

	err := db.QueryRow(query, id).Scan(&banner.ID, &banner.DeskBannerUrl, &banner.MobBannerUrl, &banner.Title, &banner.Subtitle, &banner.CtaText, &banner.CtaAction)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Banner not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(banner); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newBanner Banner
	err := json.NewDecoder(r.Body).Decode(&newBanner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO banners (deskbannerurl, mobbannerurl, title, subtitle, ctatext, ctaaction)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// Execute the insert query and get the new ID
	err = db.QueryRow(query, newBanner.DeskBannerUrl, newBanner.MobBannerUrl, newBanner.Title, newBanner.Subtitle, newBanner.CtaText, newBanner.CtaAction).Scan(&newBanner.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	banners = append(banners, newBanner)
	json.NewEncoder(w).Encode(newBanner)
}

func updateBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedBannner Banner
	err := json.NewDecoder(r.Body).Decode(&updatedBannner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedBannner.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	query := `UPDATE banners SET deskbannerurl = $1, mobbannerurl= $2,  title= $3, subtitle= $4, ctatext= $5, ctaaction = $6 WHERE id = $7;`
	result, err := db.Exec(query, updatedBannner.DeskBannerUrl, updatedBannner.MobBannerUrl, updatedBannner.Title, updatedBannner.Subtitle, updatedBannner.CtaText, updatedBannner.CtaAction, updatedBannner.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedBannner)
}

func deleteBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bannerID := r.URL.Query().Get("id")

	query := `DELETE FROM banners WHERE id = $1;`
	_, err := db.Exec(query, bannerID)

	if err != nil {
		http.Error(w, "Failed to delete banner", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Banner deleted successfully"})
}
