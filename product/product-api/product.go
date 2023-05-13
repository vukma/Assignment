package productapi

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Product model
type Product struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
}

// ProductAPI model
type ProductAPI struct {
	db  *sql.DB
	log *log.Logger
}

// New() returns a new instance of ProductAPI
func New(db *sql.DB, log *log.Logger) *ProductAPI {
	return &ProductAPI{db, log}
}

// Logger() middleware that logs info about each request coming to ProductAPI
func (p *ProductAPI) Logger(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		p.log.Printf("HTTP %s request sent to %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		n(w, r, ps)
	}
}

// Migrate() creates a postgres table called product
func (p *ProductAPI) Migrate() error {

	q := `CREATE TABLE IF NOT EXISTS product(
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(60) NOT NULL,
	description TEXT NOT NULL,
	price NUMERIC NOT NULL)`

	if _, err := p.db.Exec(q); err != nil {

		p.log.Printf("Failed to run migrations: %s\n", err)
		return err
	}

	p.log.Println("Successfully ran migrations")

	return nil
}
