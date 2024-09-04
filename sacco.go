package main

import (
	"database/sql"
	"fmt"
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
	if r.Method == http.MethodGet {
		saccos, err := GetAllSaccos(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "sacco.html", saccos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		sacco := Sacco{
			SaccoName: r.FormValue("sacco_name"),
			Manager:   r.FormValue("manager"),
			Contact:   r.FormValue("contact"),
		}
		if err := CreateSacco(db, sacco); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error adding sacco")
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// editSaccoHandler handles editing a sacco
func editSaccoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		saccoID, _ := strconv.Atoi(r.FormValue("id"))
		sacco := Sacco{
			ID:        saccoID,
			SaccoName: r.FormValue("sacco_name"),
			Manager:   r.FormValue("manager"),
			Contact:   r.FormValue("contact"),
		}
		if err := updateSacco(db, sacco); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error editing sacco")
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// deleteSaccoHandler handles deleting a sacco
func deleteSaccoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		saccoID, _ := strconv.Atoi(r.FormValue("id"))
		if err := deleteSacco(db, saccoID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error deleting sacco")
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Update sacco in the database
func updateSacco(db *sql.DB, sacco Sacco) error {
	_, err := db.Exec("UPDATE saccos SET sacco_name=?, manager=?, contact=? WHERE id=?", sacco.SaccoName, sacco.Manager, sacco.Contact, sacco.ID)
	return err
}

// Delete sacco from the database
func deleteSacco(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM saccos WHERE id=?", id)
	return err
}
