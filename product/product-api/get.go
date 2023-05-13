package productapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Healthz() writes the status of the API service to response writer
func (p *ProductAPI) Healthz(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	h := map[string]string{"status": "Service is healthy,ok"}

	if err := WriteJSON(w, http.StatusOK, h); err != nil {

		msg := fmt.Sprintf("Failed to write JSON response: %s\n", err.Error())

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}
}

// GetAll() writes all products from postgres to response writer
func (p *ProductAPI) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var products []Product

	q := `SELECT * FROM product;`

	rows, err := p.db.Query(q)

	if err != nil {

		msg := fmt.Sprintf("Failed to query Postgres: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var product Product

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)

		if err != nil {

			msg := fmt.Sprintf("Failed to scan rows: %s\n", err)

			p.log.Println("Error: ", msg)

			WriteERR(w, http.StatusBadRequest, msg)
			return
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {

		msg := fmt.Sprintf("An error occured: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	if err := WriteJSON(w, http.StatusOK, &products); err != nil {

		msg := fmt.Sprintf("Failed to write JSON response: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}
}

// GetByID writes one product from postgres to response writer
func (p *ProductAPI) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	idStr := ps.ByName("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {

		msg := fmt.Sprintf("Invalid ID provided: %s\n", err)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	var product Product

	q := `SELECT * FROM product WHERE id=$1`

	row := p.db.QueryRow(q, &id)

	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {

		msg := fmt.Sprintf("Failed to scan Product: %s\n", err)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}

	if err := WriteJSON(w, http.StatusOK, &product); err != nil {

		msg := fmt.Sprintf("Failed to write JSON response: %s\n", err)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}
}
