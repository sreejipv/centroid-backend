package main

import (
	"encoding/json"
	"net/http"

	// Import the cors package

	_ "github.com/lib/pq"
)

type Banner struct {
	ID            int    `json:"id"`
	DeskBannerUrl string `json:"deskbannerUrl"`
	MobBannerUrl  string `json:"mobbannerUrl"`
	TestData      string `json:"testData"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	CtaText       string `json:"ctaText"`
	CtaAction     string `json:"ctaAction"`
}

var banners []Banner

func getBanners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, deskbannerurl, mobbannerurl, testdata, title, subtitle, ctatext, ctaaction FROM banners")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var banners []Banner
	for rows.Next() {
		var banner Banner
		err := rows.Scan(&banner.ID, &banner.DeskBannerUrl, &banner.MobBannerUrl, &banner.TestData, &banner.Title, &banner.Subtitle, &banner.CtaText, &banner.CtaAction)
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
func createBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newBanner Banner
	err := json.NewDecoder(r.Body).Decode(&newBanner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO banners (deskbannerurl, mobbannerurl, testdata, title, subtitle, ctatext, ctaaction)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	// Execute the insert query and get the new ID
	err = db.QueryRow(query, newBanner.DeskBannerUrl, newBanner.MobBannerUrl, newBanner.TestData, newBanner.Title, newBanner.Subtitle, newBanner.CtaText, newBanner.CtaAction).Scan(&newBanner.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newBanner.ID = len(banners) + 1
	banners = append(banners, newBanner)
	json.NewEncoder(w).Encode(newBanner)

}
