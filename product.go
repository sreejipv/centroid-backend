package main

import (
	"encoding/json"
	"net/http"

	// Import the cors package

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ImgURL      string `json:"imgurl,omitempty"` // Use pointers for nullable fields
	Description string `json:"description,omitempty"`
}
type Product struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Category       []string        `json:"category,omitempty"`
	Catalog        string          `json:"catalog,omitempty"`
	FeatureDesc    string          `json:"feature_desc,omitempty"`
	FeatureList    json.RawMessage `json:"feature_list,omitempty"`
	Specifications json.RawMessage `json:"specifications,omitempty"`
	TechInfo       json.RawMessage `json:"techinfo,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	ImageGallery   []string        `json:"image_gallery,omitempty"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name, category, catalog, feature_desc, feature_list, specifications, techinfo, tags, image_gallery FROM products")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Catalog, &product.Category, &product.FeatureDesc, &product.FeatureList, &product.Specifications, &product.TechInfo, &product.Tags, &product.ImageGallery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO products (name, category, catalog, feature_desc, feature_list, specifications, techinfo, tags, image_gallery)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id;
    `
	// Execute the query
	err = db.QueryRow(query,
		newProduct.Name,
		pq.Array(newProduct.Category),
		newProduct.Catalog,
		newProduct.FeatureDesc,
		newProduct.FeatureList,
		newProduct.Specifications,
		newProduct.TechInfo,
		pq.Array(newProduct.Tags),
		pq.Array(newProduct.ImageGallery),
	).Scan(&newProduct.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Product already exists", http.StatusConflict) // HTTP 409 Conflict
				return
			}
		}
	}

	if newProduct.Name == "" {
		http.Error(w, "Empty Name is not allowed", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newProduct)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if updatedProduct.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `
        UPDATE products SET name = $1, category = $2, catalog = $3, feature_desc = $4,
        feature_list = $5, specifications = $6, techinfo = $7, tags = $8, image_gallery = $9
        WHERE id = $10;`

	_, err = db.Exec(query,
		updatedProduct.Name,
		pq.Array(updatedProduct.Category),
		updatedProduct.Catalog,
		updatedProduct.FeatureDesc,
		updatedProduct.FeatureList,
		updatedProduct.Specifications,
		updatedProduct.TechInfo,
		pq.Array(updatedProduct.Tags),
		pq.Array(updatedProduct.ImageGallery),
		updatedProduct.ID,
	)

	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedProduct)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productID := r.URL.Query().Get("id")

	query := `DELETE FROM products WHERE id = $1;`
	_, err := db.Exec(query, productID)

	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}
func createTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newTag Tag
	err := json.NewDecoder(r.Body).Decode(&newTag)

	if (err) != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newTag.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO tags ( name)
	VALUES ($1)
	RETURNING id`

	err = db.QueryRow(query, newTag.Name).Scan(newTag.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Tag already exists", http.StatusConflict) // HTTP 409 Conflict
				return
			}
		}
	}

	json.NewEncoder(w).Encode(newTag)
}

func getTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name FROM tags")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(tags)
}
func updateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedTag Tag
	err := json.NewDecoder(r.Body).Decode(&updatedTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedTag.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	query := `UPDATE tags SET name = $1 WHERE id = $2;`
	result, err := db.Exec(query, updatedTag.Name, updatedTag.ID)

	if err != nil {
		http.Error(w, "Failed to update tag", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedTag)
}

func deleteTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tagID := r.URL.Query().Get("id")

	query := `DELETE FROM tags WHERE id = $1;`
	result, err := db.Exec(query, tagID)

	if err != nil {
		http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to check affected rows ", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Tag is not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Tag deleted successfully"})
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)

	if (err) != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newCategory.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO categories ( name, imgurl, description)
	VALUES ($1, $2, $3)
	RETURNING id`

	err = db.QueryRow(query, newCategory.Name, newCategory.ImgURL, newCategory.Description).Scan(&newCategory.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Tag already exists", http.StatusConflict) // HTTP 409 Conflict
				return

			}
		}
	}

	json.NewEncoder(w).Encode(newCategory)
}
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, name, imgurl, description FROM categories")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.ImgURL, &category.Description)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(categories)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedCategory Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE categories SET name = $1, imgurl = $2, description = $3 WHERE id = $4;`
	result, err := db.Exec(query, updatedCategory.Name, updatedCategory.ImgURL, updatedCategory.Description, updatedCategory.ID)

	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check the affected rows", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Category is not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updatedCategory)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categoryID := r.URL.Query().Get("id")

	query := `DELETE FROM categories WHERE id = $1;`
	result, err := db.Exec(query, categoryID)

	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check the affected rows", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Category is not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}
