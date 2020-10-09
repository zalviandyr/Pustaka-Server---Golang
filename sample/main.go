package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// Customers struct
type Customers struct {
	CustomerID   string `json:"customer_id"`
	CompanyName  string `json:"company_name"`
	ContactName  string `json:"contact_name"`
	ContactTitle string `json:"contact_title"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Region       string `json:"region"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Fax          string `json:"fax"`
}

// Response struct
type Response struct {
	StatusCode    int         `json:"status_code"`
	Message       string      `json:"message"`
	CustomersList []Customers `json:"customers_list"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response Response
		var customersList []Customers

		sql := "SELECT * FROM customers"
		result, err := db.Query(sql)

		defer func() {
			if result != nil {
				result.Close()
			}
		}()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			for result.Next() {
				var customer Customers

				result.Scan(
					&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
					&customer.ContactTitle, &customer.Address, &customer.City, &customer.Region,
					&customer.PostalCode, &customer.Country, &customer.Phone, &customer.Fax)

				customersList = append(customersList, customer)
			}

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"
			response.CustomersList = customersList
		}

		json.NewEncoder(w).Encode(response)
	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response Response
		var customersList []Customers
		params := mux.Vars(r)

		sql := "SELECT * FROM customers WHERE CustomerID = ?"
		result, err := db.Query(sql, params["id"])

		defer func() {
			if result != nil {
				result.Close()
			}
		}()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel it"
		} else {
			if result.Next() {
				var customer Customers

				result.Scan(
					&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
					&customer.ContactTitle, &customer.Address, &customer.City, &customer.Region,
					&customer.PostalCode, &customer.Country, &customer.Phone, &customer.Fax)

				customersList = append(customersList, customer)
			}

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"
			response.CustomersList = customersList
		}

		json.NewEncoder(w).Encode(response)
	}
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var response Response

		customerID := r.FormValue("customer_id")
		companyName := r.FormValue("company_name")
		contactName := r.FormValue("contact_name")
		contactTitle := r.FormValue("contact_title")
		address := r.FormValue("address")
		city := r.FormValue("city")
		region := r.FormValue("region")
		postalCode := r.FormValue("postal_code")
		country := r.FormValue("country")
		phone := r.FormValue("phone")
		fax := r.FormValue("fax")

		sql := `INSERT INTO customers VALUES(?,?,?, ?,?,? ,?,?,? ,?,?)`
		stmt, err := db.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			_, err := stmt.Exec(customerID, companyName, contactName, contactTitle,
				address, city, region, postalCode, country, phone, fax)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data duplicated"
			} else {
				w.WriteHeader(http.StatusCreated)
				response.StatusCode = http.StatusCreated
				response.Message = "Data created"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var response Response
		params := mux.Vars(r)

		companyName := r.FormValue("company_name")
		contactName := r.FormValue("contact_name")
		contactTitle := r.FormValue("contact_title")
		address := r.FormValue("address")
		city := r.FormValue("city")
		region := r.FormValue("region")
		postalCode := r.FormValue("postal_code")
		country := r.FormValue("country")
		phone := r.FormValue("phone")
		fax := r.FormValue("fax")

		sql := `UPDATE customers SET 
		CompanyName = ?, ContactName = ?, ContactTitle = ?, Address = ?, City = ?,
		Region = ?, PostalCode = ?, Country = ?, Phone = ?, Fax = ? 
		WHERE CustomerID = ?`

		stmt, err := db.Prepare(sql)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			result, _ := stmt.Exec(companyName, contactName, contactTitle, address, city,
				region, postalCode, country, phone, fax, params["id"])
			rowAffect, _ := result.RowsAffected()

			if rowAffect == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data not found"
			} else {
				w.WriteHeader(http.StatusCreated)
				response.StatusCode = http.StatusCreated
				response.Message = "Data updated"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "DELETE" {
		var response Response
		params := mux.Vars(r)

		sql := "DELETE FROM customers WHERE CustomerID = ?"
		stmt, err := db.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			result, _ := stmt.Exec(params["id"])
			rowAffected, _ := result.RowsAffected()

			if rowAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Data deleted"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// init database
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// init router
	router := mux.NewRouter()

	// routing
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", createCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
