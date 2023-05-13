# Product API

## Overview:
ProductAPI is a REST API that enables basic CRUD operations to manage the state of Products. \
The API utilizes several packages, including `godotenv` for managing environment variables, `httprouter` package for managing HTTP server 
routing, `validator/v10` for validating `json` objects and `lib/pq` as PostgreSQL driver.

## Running the project locally:
To run the ProductAPI locally, you can use either the `go` binary or the `make` command, which runs commands based on the Makefile provided in the repository.

Run with `go`:

```bash
# Build the project and store the binary in ./bin/ directory 
go build -o ./bin/product *.go

# Run the built binary
./bin/product 
```

Run with `make`:
```bash
make run
```


## Models:
**Product model:**
```json
{
"id": <int>,
"name": <string>,
"description": <string>,
"price": <float32>,
}
```

**Product request model:**
```json
{
"name": <string>,
"description": <string>,
"price": <float32>,
}
```

**Health model**
```json
{
"status":<string>
}
```

**Error model**
```json
{
"error":<string>
}
```

## ProductAPI Endpoints:
ProductAPI has 6 endpoints:
- `GET /product`: Retrieves all existing products from the database. \
  Example output:
  ```json
  [
    {
    "id": 1, 
    "name": "Cheeseburger", 
    "description": "Juicy cheeseburger made with fresh ingredients",
    "price": 5.78,
    },
    {
    "id": 2,
    "name": "Coffee",
    "description": "Strong black coffee without sugar",
    "price": 1.22,
    }
  ]
  ```

- `GET /product/1`: Retrieves one product from the database based on the provided ID. \
  Example output:
  ```json
    {
    "id": 1, 
    "name": "Cheeseburger", 
    "description": "Juicy cheeseburger made with fresh ingredients",
    "price": 5.78,
    }
  ```

- `GET /healthz`: Retrieves the status of the API. \
  Example output:
  ```json
    {
    "status": "Service is healthy,ok"
    }
  ```

- `POST /product`: Posts a new product and adds it to the database. \
  Example payload:
  ```json
    {
    "name": "Green Tea",
    "description": "Healthy green tea with a little bit of honey",
    "price": 1.00, 
    }
  ```
  Example output:

  ```json
    {
    "insertedID": 3
    }
  ```

- `PUT /product/3`: Updates existing product. \
  Example payload:
  ```json
    {
    "name": "Vegan salad",
    "description": "Fruit and vegetable salad with tofu",
    "price": 4.20, 
    }
  ```
  Example output:

  ```json
    {
    "id": 3,
    "name": "Vegan salad",
    "description": "Fruit and vegetable salad with tofu",
    "price": 4.20, 
    }
  ```

- `DELETE /product/3`: Deletes a product from the database. \
  On success it returns a `204` (No Content) header without any output
