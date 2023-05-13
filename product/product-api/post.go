package productapi

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// PostOne() adds new product to postgres and writes it to response writer
func (p *ProductAPI) PostOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var product Product

	if err := ReadJSON(r, &product); err != nil {

		msg := fmt.Sprintf("Failed to read JSON body: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	if err := ValidateJSON(&product); err != nil {

		msg := fmt.Sprintf("Failed to validate JSON: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	q := `INSERT INTO product(name, description, price)
	VALUES ($1,$2,$3)
	RETURNING id`

	var id int

	err := p.db.QueryRow(q, &product.Name, &product.Description, &product.Price).Scan(&id)

	if err != nil {

		msg := fmt.Sprintf("Failed to query Postgres: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	if err := WriteJSON(w, http.StatusCreated, map[string]int{"insertedID": id}); err != nil {

		msg := fmt.Sprintf("Failed to write JSON response: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}
}
