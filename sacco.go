package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

type Sacco struct {
	ID        int
	SaccoName string
	Manager   string
	Contact   string
}

// GetAllSaccos retrieves all saccos from the database
func GetAllSaccos(db *sql.DB) ([]Sacco, error) {
	rows, err := db.Query("SELECT id, sacco_name, manager, contact FROM saccos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saccos []Sacco
	for rows.Next() {
		var sacco Sacco
		if err := rows.Scan(&sacco.ID, &sacco.SaccoName, &sacco.Manager, &sacco.Contact); err != nil {
			return nil, err
		}
		saccos = append(saccos, sacco)
	}

	return saccos, nil
}

// CreateSacco inserts a new Sacco into the database
func CreateSacco(db *sql.DB, sacco Sacco) error {
	_, err := db.Exec("INSERT INTO saccos (sacco_name, manager, contact) VALUES (?, ?, ?)", sacco.SaccoName, sacco.Manager, sacco.Contact)
	return err
}

// saccoHandler handles requests to the /home route
func saccoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		saccos, err := GetAllSaccos(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "sacco.html", saccos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		switch r.URL.Path {
		case "/home":
			sacco := Sacco{
				SaccoName: r.FormValue("sacco_name"),
				Manager:   r.FormValue("manager"),
				Contact:   r.FormValue("contact"),
			}
			if err := CreateSacco(db, sacco); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	case "/edit-sacco":
		if r.URL.Path == "/edit-sacco" {
			// Handle delete sacco
			saccoID, _ := strconv.Atoi(r.FormValue("id"))
			sacco := Sacco{
				ID:        saccoID,
				SaccoName: r.FormValue("sacco_name"),
				Manager:   r.FormValue("manager"),
				Contact:   r.FormValue("contact"),
			}
			if err := updateSacco(db, sacco); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.URL.Path == "/delete-sacco" {
			// Handle delete sacco
			saccoID, _ := strconv.Atoi(r.FormValue("id"))
			if err := deleteSacco(db, saccoID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Update sacco in the database
func updateSacco(db *sql.DB, sacco Sacco) error {
	_, err := db.Exec("Update saccos SET sacco_name=?, manager=?, contact=? WHERE id=?", sacco.SaccoName, sacco.Manager, sacco.ID)
	return err
}

// Delete sacco from the database
func deleteSacco(db *sql.DB, id int) error {
	_, err := db.Exec("Delete FROM saccos WHERE id=?", id)
	return err
}
