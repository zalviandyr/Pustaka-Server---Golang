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

// Buku struct
type Buku struct {
	ISBN        string `json:"isbn"`
	Judul       string `json:"judul"`
	Penulis     string `json:"penulis"`
	Penerbit    string `json:"penerbit"`
	TahunTerbit string `json:"tahun_terbit"`
	Cover       string `json:"cover"`
}

// Response struct
type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	BukuList   []Buku `json:"buku_list"`
}

func index(w http.ResponseWriter, r *http.Request) {
	var message = "Welcome to Pustaka Server"
	w.Write([]byte(message))
}

func getAllBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response Response
		var bukuList []Buku

		sql := `SELECT * FROM buku`
		result, err := db.Query(sql)

		// defer berfungsi untuk mengeksekusi suatu fungsi
		// di bagian akhir dari fungsi tersebut
		// misal ny fungsi result.Close()
		// akan dieksekusi pada bagian akhir dari fungsi getAllBuku()

		defer func() {
			// jika terjadi error maka result akan null
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
				var buku Buku
				result.Scan(&buku.ISBN, &buku.Judul,
					&buku.Penulis, &buku.Penerbit,
					&buku.TahunTerbit, &buku.Cover)

				bukuList = append(bukuList, buku)
			}

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"
			response.BukuList = bukuList
		}

		resultJSON, _ := json.Marshal(response)
		w.Write(resultJSON)
	}
}

func getBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response Response
		var bukuList []Buku

		params := mux.Vars(r)
		sql := `SELECT * FROM buku WHERE isbn = ?`
		result, err := db.Query(sql, params["id"])

		// defer berfungsi untuk mengeksekusi suatu fungsi
		// di bagian akhir dari fungsi tersebut
		// misal ny fungsi result.Close()
		// akan dieksekusi pada bagian akhir dari fungsi getAllBuku()

		defer func() {
			// jika terjadi error maka result akan null
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
				var buku Buku
				result.Scan(&buku.ISBN, &buku.Judul,
					&buku.Penulis, &buku.Penerbit,
					&buku.TahunTerbit, &buku.Cover)

				bukuList = append(bukuList, buku)
			}

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"
			response.BukuList = bukuList
		}

		resultJSON, _ := json.Marshal(response)
		w.Write(resultJSON)
	}
}

func createBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var response Response

		isbn := r.FormValue("isbn")
		judul := r.FormValue("judul")
		penulis := r.FormValue("penulis")
		penerbit := r.FormValue("penerbit")
		tahunTerbit := r.FormValue("tahun_terbit")
		cover := r.FormValue("cover")

		stmt, err := db.Prepare("INSERT INTO buku (isbn, judul, penulis, penerbit, tahun_terbit, cover) values (?,?,?,?,?,?)")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			_, err := stmt.Exec(isbn, judul, penulis, penerbit, tahunTerbit, cover)

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

		resultJSON, _ := json.Marshal(response)
		w.Write(resultJSON)
	}
}

func updateBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var response Response
		params := mux.Vars(r)

		judul := r.FormValue("judul")
		penulis := r.FormValue("penulis")
		penerbit := r.FormValue("penerbit")
		tahunTerbit := r.FormValue("tahun_terbit")
		cover := r.FormValue("cover")

		stmt, err := db.Prepare(`UPDATE buku SET 
		judul = ?, 
		penulis = ?,
		penerbit = ?,
		tahun_terbit = ?,
		cover = ?
		WHERE isbn = ?`)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			result, _ := stmt.Exec(judul, penulis, penerbit, tahunTerbit, cover, params["id"])
			rowAffect, _ := result.RowsAffected()

			if rowAffect == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data not found or has been updated"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Data success update"
			}
		}

		resultJSON, _ := json.Marshal(response)
		w.Write(resultJSON)
	}
}

func deleteBuku(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "DELETE" {
		var response Response
		params := mux.Vars(r)

		stmt, err := db.Prepare("DELETE FROM buku WHERE isbn = ?")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Something wrong, i can feel"
		} else {
			result, _ := stmt.Exec(params["id"])
			rowAffect, _ := result.RowsAffected()

			if rowAffect == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Data success delete"
			}
		}

		resultJSON, _ := json.Marshal(response)
		w.Write(resultJSON)
	}
}

func main() {
	// init database
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/pustaka")
	if err != nil {
		// panic func adalah untuk membuat program keluar langsung
		panic(err.Error())
	}

	defer db.Close()

	// init router
	router := mux.NewRouter()

	// router handler & endpoints
	router.HandleFunc("/", index)
	router.HandleFunc("/buku", getAllBuku).Methods("GET")
	router.HandleFunc("/buku/{id}", getBuku).Methods("GET")
	router.HandleFunc("/buku/{id}", updateBuku).Methods("PUT")
	router.HandleFunc("/buku/{id}", deleteBuku).Methods("DELETE")
	router.HandleFunc("/buku", createBuku).Methods("POST")

	// start server
	log.Fatal(http.ListenAndServe(":8080", router))
}
