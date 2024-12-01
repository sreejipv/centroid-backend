package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// Import the cors package
	"github.com/lib/pq"
)

type Order struct {
	ID               string `json:"id"`
	Invoice_number   string `json:"invoice_number"`
	Transporter_name string `json:"transporter_name"`
	LR_number        string `json:"lr_number"`
	Batch_number     string `json:"batch_number"`
	Remarks          string `json:"remarks"`
	Order_status     string `json:"order_status"`
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)

	if (err) != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newOrder.Invoice_number == "" {
		http.Error(w, "Order Number ID be empty", http.StatusBadRequest)
		return
	}
	query := `INSERT INTO orders ( invoice_number, transporter_name, lr_number, batch_number, remarks, order_status  )
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`

	err = db.QueryRow(query, newOrder.Invoice_number, newOrder.Transporter_name, newOrder.LR_number, newOrder.Batch_number, newOrder.Remarks, newOrder.Order_status).Scan(newOrder.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				http.Error(w, "Order already exists", http.StatusConflict) // HTTP 409 Conflict
				return
			}
		}
	}

	json.NewEncoder(w).Encode(newOrder)
}

func fetchOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing Order ID", http.StatusBadRequest)
		return
	}

	query := `SELECT  id, invoice_number, transporter_name,lr_number, batch_number, remarks, order_status  FROM orders WHERE id=$1`
	var order Order

	err := db.QueryRow(query, id).Scan(&order.ID, &order.Invoice_number, &order.Transporter_name, &order.LR_number, &order.Batch_number, &order.Remarks, &order.Order_status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	rows, err := db.Query("SELECT id, invoice_number, transporter_name,lr_number, batch_number, remarks, order_status FROM orders")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Invoice_number, &order.Transporter_name, &order.LR_number, &order.Batch_number, &order.Remarks, &order.Order_status)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(orders)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedOrder Order
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedOrder.Invoice_number == "" {
		http.Error(w, "Invoice number cannot be empty", http.StatusBadRequest)
		return
	}

	query := `UPDATE orders SET invoice_number  = $1, transporter_name  = $2, lr_number  = $3, batch_number = $4, remarks = $5, order_status = $6 WHERE id = $7;`
	result, err := db.Exec(query, updatedOrder.Invoice_number, updatedOrder.Transporter_name, updatedOrder.LR_number, updatedOrder.Batch_number, updatedOrder.Remarks, updatedOrder.Order_status, updatedOrder.ID)

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
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderID := r.URL.Query().Get("id")

	query := `DELETE FROM orders WHERE id = $1;`
	result, err := db.Exec(query, orderID)

	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
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

	json.NewEncoder(w).Encode(map[string]string{"message": "order deleted successfully"})
}
