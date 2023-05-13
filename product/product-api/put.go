package productapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// UpdateOne() updates one product in postgres and writes it to response writer
func (p *ProductAPI) UpdateOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	idStr := ps.ByName("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {

		msg := fmt.Sprintf("Invalid ID provided: %s", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	// id := int64(idInt)

	var product Product

	if err := ReadJSON(r, &product); err != nil {

		msg := fmt.Sprintf("Failed to read JSON body: %s", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	if err := ValidateJSON(&product); err != nil {

		msg := fmt.Sprintf("Failed to validate JSON: %s", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	q := `UPDATE product
	SET name=$2, description=$3, price=$4
	WHERE id=$1
	RETURNING *`

	row := p.db.QueryRow(q,
		&id,
		product.Name,
		product.Description,
		product.Price)

	err = row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
	)

	if err != nil {

		msg := fmt.Sprintf("Failed to scan Product: %s", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}

	if err := WriteJSON(w, http.StatusOK, &product); err != nil {

		msg := fmt.Sprintf("Failed to write JSON response: %s", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusInternalServerError, msg)
		return
	}
}
