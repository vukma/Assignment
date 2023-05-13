package productapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// DeleteOne() deletes one product from postgres
func (p *ProductAPI) DeleteOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	idStr := ps.ByName("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {

		msg := fmt.Sprintf("Invalid ID provided: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	q := `DELETE FROM product WHERE id=$1`

	result, err := p.db.Exec(q, &id)

	if err != nil {

		msg := fmt.Sprintf("Failed to query Postgres: %s\n", err)

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {

		msg := fmt.Sprintf("Failed to query affected rows: %s\n", err)

		WriteERR(w, http.StatusBadRequest, msg)
		return
	}

	if rowsAffected == 0 {

		msg := "Couldn't delete resource: Resource not found"

		p.log.Println("Error: ", msg)

		WriteERR(w, http.StatusNotFound, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
